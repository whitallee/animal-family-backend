package enclosure

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// ErrAnimalNotOwned is returned by CreateEnclosureWithAnimals when animalIds
// names an animal the caller does not own. Handlers map it to 403 so it is not
// reported as a server fault.
var ErrAnimalNotOwned = errors.New("animal does not exist or does not belong to you")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateEnclosure(enclosure types.Enclosure, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var addedEnclosureId int
	err = tx.QueryRow(`INSERT INTO "enclosures" ("enclosureName", "image", "notes", "habitatId") VALUES ($1,$2,$3,$4) RETURNING "enclosureId"`, enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId).Scan(&addedEnclosureId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO "enclosureUser" ("enclosureId", "userId") VALUES ($1,$2)`, addedEnclosureId, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateEnclosureWithAnimals(enclosure types.Enclosure, animalIds []int, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	// Rolls back every early return below. A no-op once Commit has succeeded.
	defer func() { _ = tx.Rollback() }()

	var addedEnclosureId int
	err = tx.QueryRow(`INSERT INTO "enclosures" ("enclosureName", "image", "notes", "habitatId") VALUES ($1,$2,$3,$4) RETURNING "enclosureId"`, enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId).Scan(&addedEnclosureId)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO "enclosureUser" ("enclosureId", "userId") VALUES ($1,$2)`, addedEnclosureId, userID); err != nil {
		return err
	}

	if err := assignAnimalsToEnclosure(tx, addedEnclosureId, animalIds, userID); err != nil {
		return err
	}

	return tx.Commit()
}

// assignAnimalsToEnclosure moves the caller's animals into the enclosure.
//
// The EXISTS clause is what stops a caller moving somebody else's animal into
// their own enclosure. Without it any animalId was accepted, which both exposed
// the victim's animal through the enclosure's animal list and let it be
// destroyed via a cascading enclosure delete.
//
// The check lives inside the UPDATE rather than in a preceding SELECT so that
// it is atomic: under READ COMMITTED another session could otherwise change
// ownership between a separate check and the write.
//
// One statement covers every animal. The previous loop issued a round trip per
// animal and stopped at the first bad one, so a caller correcting a bulk
// request learned about their mistakes one at a time.
func assignAnimalsToEnclosure(tx *sql.Tx, enclosureId int, animalIds []int, userID int) error {
	wanted := dedupe(animalIds)
	if len(wanted) == 0 {
		return nil
	}

	result, err := tx.Exec(
		`UPDATE "animals" SET "enclosureId" = $1
		 WHERE "animalId" = ANY($2)
		   AND EXISTS (SELECT 1 FROM "animalUser"
		               WHERE "animalUser"."animalId" = "animals"."animalId"
		                 AND "animalUser"."userId" = $3)`,
		enclosureId, pq.Array(wanted), userID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Deduplicating above is what makes this comparison meaningful: it counts
	// distinct rows, so it no longer depends on how a driver reports an update
	// that writes a value a row already holds.
	if affected == int64(len(wanted)) {
		return nil
	}

	// Only reached when something is already wrong, so the extra query costs
	// nothing on the happy path and buys an error naming every offending id.
	unowned, err := findUnownedAnimals(tx, wanted, userID)
	if err != nil {
		return err
	}

	return fmt.Errorf("animals %v: %w", unowned, ErrAnimalNotOwned)
}

// findUnownedAnimals returns the requested ids that do not belong to the user,
// including ids that do not exist at all.
func findUnownedAnimals(tx *sql.Tx, wanted []int, userID int) ([]int, error) {
	rows, err := tx.Query(
		`SELECT "animalId" FROM "animalUser" WHERE "userId" = $1 AND "animalId" = ANY($2)`,
		userID, pq.Array(wanted),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	owned := make(map[int]bool, len(wanted))
	for rows.Next() {
		var animalId int
		if err := rows.Scan(&animalId); err != nil {
			return nil, err
		}
		owned[animalId] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	unowned := make([]int, 0, len(wanted))
	for _, animalId := range wanted {
		if !owned[animalId] {
			unowned = append(unowned, animalId)
		}
	}

	return unowned, nil
}

// dedupe preserves order so error messages list ids as the caller sent them.
func dedupe(ids []int) []int {
	seen := make(map[int]bool, len(ids))
	unique := make([]int, 0, len(ids))

	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		unique = append(unique, id)
	}

	return unique
}

func (s *Store) UpdateEnclosure(enclosure types.Enclosure) error {
	_, err := s.db.Exec(`UPDATE "enclosures"
						SET "enclosureName" = $1, "image" = $2, "notes" = $3, "habitatId" = $4
						WHERE "enclosureId" = $5`, enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId, enclosure.EnclosureId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateEnclosureOwnerWithAnimals(oldEnclosureUser types.EnclosureUser, newUserId int) error {
	// get animals in enclosure
	rows, err := s.db.Query(`SELECT "animalId" FROM "animals" WHERE "enclosureId" = $1`, oldEnclosureUser.EnclosureId)
	if err != nil {
		return err
	}
	animalIds := make([]int, 0)
	for rows.Next() {
		var animalId int
		if err := rows.Scan(&animalId); err != nil {
			_ = rows.Close()
			return err
		}
		animalIds = append(animalIds, animalId)
	}
	_ = rows.Close()

	// start update owners transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// update all animals' owners
	for _, animalId := range animalIds {
		_, err := tx.Exec(`UPDATE "animalUser"
							SET "userId" = $1
							WHERE "animalId" = $2 AND "userId" = $3`, newUserId, animalId, oldEnclosureUser.UserID)
		if err != nil {
			return err
		}
	}

	// update enclosure owner
	_, err = tx.Exec(`UPDATE "enclosureUser"
						SET "userId" = $1
						WHERE "enclosureId" = $2 AND "userId" = $3`, newUserId, oldEnclosureUser.EnclosureId, oldEnclosureUser.UserID)
	if err != nil {
		return err
	}

	// commit changes
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetEnclosures() ([]*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT * FROM "enclosures"`)
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		s, err := utils.ScanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}

		enclosures = append(enclosures, s)
	}

	return enclosures, nil
}

func (s *Store) GetEnclosureByNameAndHabitatWithUserId(enclosureName string, habitatId int, userID int) (*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT e."enclosureId", e."enclosureName", e."image", e."Notes", e."habitatId"
							FROM "enclosures" e JOIN "enclosureUser" ON "enclosureUser"."enclosureId"=e."enclosureId"
							WHERE "enclosureName" = $1 AND "habitatId" = $2 AND "userId" = $3`, enclosureName, habitatId, userID)
	if err != nil {
		return nil, err
	}

	enclosure := new(types.Enclosure)
	for rows.Next() {
		enclosure, err = utils.ScanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}
	}

	if enclosure.EnclosureId == 0 {
		return nil, fmt.Errorf("enclosure not found")
	}

	return enclosure, nil
}

func (s *Store) GetEnclosureUserByIds(enclosureId int, userID int) (*types.EnclosureUser, error) {
	rows, err := s.db.Query(`SELECT * FROM "enclosureUser" WHERE "enclosureId" = $1 AND "userId" = $2`, enclosureId, userID)
	if err != nil {
		return nil, err
	}

	enclosureUser := new(types.EnclosureUser)
	for rows.Next() {
		enclosureUser, err = utils.ScanRowsIntoEnclosureUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if enclosureUser.EnclosureId == 0 && enclosureUser.UserID == 0 {
		return nil, fmt.Errorf("no ownership found between user and enclosure")
	}

	return enclosureUser, nil
}

func (s *Store) GetEnclosureUserByEnclosureId(enclosureId int) (*types.EnclosureUser, error) {
	rows, err := s.db.Query(`SELECT * FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return nil, err
	}

	enclosureUser := new(types.EnclosureUser)
	for rows.Next() {
		enclosureUser, err = utils.ScanRowsIntoEnclosureUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if enclosureUser.EnclosureId == 0 && enclosureUser.UserID == 0 {
		return nil, fmt.Errorf("no owner of enclosure with id %d found", enclosureId)
	}

	return enclosureUser, nil
}

func (s *Store) GetEnclosuresByUserId(userID int) ([]*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT e."enclosureId", e."enclosureName", e."image", e."notes", e."habitatId"
							FROM "enclosures" e JOIN "enclosureUser" ON "enclosureUser"."enclosureId"=e."enclosureId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		enclosure, err := utils.ScanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}

		enclosures = append(enclosures, enclosure)
	}

	return enclosures, nil
}

func (s *Store) GetEnclosureById(enclosureId int) (*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT * FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		enclosure, err := utils.ScanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}

		enclosures = append(enclosures, enclosure)
	}

	if len(enclosures) == 0 {
		return nil, nil
	}

	return enclosures[0], nil
}

func (s *Store) DeleteEnclosureById(enclosureId int) error {
	// get tasks associated with this enclosure
	taskRows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	taskIds := make([]int, 0)
	for taskRows.Next() {
		var taskId int
		err = taskRows.Scan(&taskId)
		if err != nil {
			_ = taskRows.Close()
			return err
		}
		taskIds = append(taskIds, taskId)
	}
	_ = taskRows.Close()

	// get animals from enclosure
	animalRows, err := s.db.Query(`SELECT "animalId" FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	enclosureAnimalIds := make([]int, 0)
	for animalRows.Next() {
		var animalId int
		if err := animalRows.Scan(&animalId); err != nil {
			_ = animalRows.Close()
			return err
		}
		enclosureAnimalIds = append(enclosureAnimalIds, animalId)
	}
	_ = animalRows.Close()

	// start enclosureId updates on animals and deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete enclosure tasks first
	for _, taskId := range taskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// update enclosureId for animals
	for _, animalId := range enclosureAnimalIds {
		_, err = tx.Exec(`UPDATE "animals" SET "enclosureId" = NULL WHERE "animalId" = $1`, animalId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteEnclosureAndTasksById(enclosureId int) error {
	// get tasks associated with this enclosure
	taskRows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	taskIds := make([]int, 0)
	for taskRows.Next() {
		var taskId int
		err = taskRows.Scan(&taskId)
		if err != nil {
			_ = taskRows.Close()
			return err
		}
		taskIds = append(taskIds, taskId)
	}
	_ = taskRows.Close()

	// get animals from enclosure
	animalRows, err := s.db.Query(`SELECT "animalId" FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	enclosureAnimalIds := make([]int, 0)
	for animalRows.Next() {
		var animalId int
		if err := animalRows.Scan(&animalId); err != nil {
			_ = animalRows.Close()
			return err
		}
		enclosureAnimalIds = append(enclosureAnimalIds, animalId)
	}
	_ = animalRows.Close()

	// start deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete tasks and their related records
	for _, taskId := range taskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// update enclosureId for animals to NULL
	for _, animalId := range enclosureAnimalIds {
		_, err = tx.Exec(`UPDATE "animals" SET "enclosureId" = NULL WHERE "animalId" = $1`, animalId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteEnclosureAndAnimalsAndTasksById(enclosureId int) error {
	// get animals from enclosure
	animalRows, err := s.db.Query(`SELECT "animalId" FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	enclosureAnimalIds := make([]int, 0)
	for animalRows.Next() {
		var animalId int
		if err := animalRows.Scan(&animalId); err != nil {
			_ = animalRows.Close()
			return err
		}
		enclosureAnimalIds = append(enclosureAnimalIds, animalId)
	}
	_ = animalRows.Close()

	// get tasks associated with this enclosure
	enclosureTaskRows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	enclosureTaskIds := make([]int, 0)
	for enclosureTaskRows.Next() {
		var taskId int
		err = enclosureTaskRows.Scan(&taskId)
		if err != nil {
			_ = enclosureTaskRows.Close()
			return err
		}
		enclosureTaskIds = append(enclosureTaskIds, taskId)
	}
	_ = enclosureTaskRows.Close()

	// get all animal task IDs before starting transaction
	allAnimalTaskIds := make([]int, 0)
	for _, animalId := range enclosureAnimalIds {
		animalTaskRows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "animalId" = $1`, animalId)
		if err != nil {
			return err
		}

		for animalTaskRows.Next() {
			var taskId int
			err = animalTaskRows.Scan(&taskId)
			if err != nil {
				_ = animalTaskRows.Close()
				return err
			}
			allAnimalTaskIds = append(allAnimalTaskIds, taskId)
		}
		_ = animalTaskRows.Close()
	}

	// start deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete tasks for each animal and the animals themselves
	for _, taskId := range allAnimalTaskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// delete animals
	for _, animalId := range enclosureAnimalIds {
		_, err = tx.Exec(`DELETE FROM "animalUser" WHERE "animalId" = $1`, animalId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "animals" WHERE "animalId" = $1`, animalId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// delete enclosure tasks
	for _, taskId := range enclosureTaskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UserOwnsEnclosure reports whether the user owns the enclosure. See the note
// on animal.Store.UserOwnsAnimal for why this exists alongside
// GetEnclosureUserByIds.
func (s *Store) UserOwnsEnclosure(enclosureId int, userID int) (bool, error) {
	var owned bool
	err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM "enclosureUser" WHERE "enclosureId" = $1 AND "userId" = $2)`,
		enclosureId, userID,
	).Scan(&owned)
	if err != nil {
		return false, err
	}

	return owned, nil
}
