package enclosure

import (
	"database/sql"
	"fmt"

	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

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

	var addedEnclosureId int
	err = tx.QueryRow(`INSERT INTO "enclosures" ("enclosureName", "image", "notes", "habitatId") VALUES ($1,$2,$3,$4) RETURNING "enclosureId"`, enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId).Scan(&addedEnclosureId)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO "enclosureUser" ("enclosureId", "userId") VALUES ($1,$2)`, addedEnclosureId, userID); err != nil {
		return err
	}

	for _, animalId := range animalIds {
		if _, err := tx.Exec(`UPDATE "animals" SET "enclosureId" = $1 WHERE "animalId" = $2`, addedEnclosureId, animalId); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
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
	rows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc"
							FROM "animals" WHERE "enclosureId" = $1`, oldEnclosureUser.EnclosureId)
	if err != nil {
		return err
	}
	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return err
		}

		animals = append(animals, animal)
	}

	// start update owners transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// update all animals' owners
	for _, animal := range animals {
		_, err := tx.Exec(`UPDATE "animalUser"
							SET "userId" = $1
							WHERE "animalId" = $2 AND "userId" = $3`, newUserId, animal.AnimalId, oldEnclosureUser.UserID)
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
			taskRows.Close()
			return err
		}
		taskIds = append(taskIds, taskId)
	}
	taskRows.Close()

	// get animals from enclosure
	animalRows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc"
							FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	animals := make([]*types.Animal, 0)
	for animalRows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(animalRows)
		if err != nil {
			animalRows.Close()
			return err
		}
		animals = append(animals, animal)
	}
	animalRows.Close()

	// start enclosureId updates on animals and deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete enclosure tasks first
	for _, taskId := range taskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// update enclosureId for animals
	for _, animal := range animals {
		_, err = tx.Exec(`UPDATE "animals" SET "enclosureId" = NULL WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
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
			taskRows.Close()
			return err
		}
		taskIds = append(taskIds, taskId)
	}
	taskRows.Close()

	// get animals from enclosure
	animalRows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc"
							FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	animals := make([]*types.Animal, 0)
	for animalRows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(animalRows)
		if err != nil {
			animalRows.Close()
			return err
		}
		animals = append(animals, animal)
	}
	animalRows.Close()

	// start deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete tasks and their related records
	for _, taskId := range taskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// update enclosureId for animals to NULL
	for _, animal := range animals {
		_, err = tx.Exec(`UPDATE "animals" SET "enclosureId" = NULL WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
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
	animalRows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc"
							FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return err
	}

	animals := make([]*types.Animal, 0)
	for animalRows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(animalRows)
		if err != nil {
			animalRows.Close()
			return err
		}
		animals = append(animals, animal)
	}
	animalRows.Close()

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
			enclosureTaskRows.Close()
			return err
		}
		enclosureTaskIds = append(enclosureTaskIds, taskId)
	}
	enclosureTaskRows.Close()

	// get all animal task IDs before starting transaction
	allAnimalTaskIds := make([]int, 0)
	for _, animal := range animals {
		animalTaskRows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			return err
		}

		for animalTaskRows.Next() {
			var taskId int
			err = animalTaskRows.Scan(&taskId)
			if err != nil {
				animalTaskRows.Close()
				return err
			}
			allAnimalTaskIds = append(allAnimalTaskIds, taskId)
		}
		animalTaskRows.Close()
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
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete animals
	for _, animal := range animals {
		_, err = tx.Exec(`DELETE FROM "animalUser" WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "animals" WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete enclosure tasks
	for _, taskId := range enclosureTaskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
