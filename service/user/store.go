package user

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

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(`INSERT INTO "users" ("firstName", "lastName", "email", "password") VALUES ($1, $2, $3, $4)`, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query(`SELECT * FROM "users" WHERE "email" = $1`, email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query(`SELECT * FROM "users" WHERE "userId" = $1`, id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) DeleteUserById(userID int) error {
	// get tasks, animals, and enclosures from userID
	tRows, err := s.db.Query(`SELECT t."taskId", t."taskName", t."complete", t."lastCompleted", t."repeatIntervHours"
							FROM "tasks" t JOIN "taskUser" ON "taskUser"."taskId"=t."taskId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}

	aRows, err := s.db.Query(`SELECT a."animalId", a."animalName", a."image", a."notes", a."speciesId", a."enclosureId"
							FROM "animals" a JOIN "animalUser" ON "animalUser"."animalId"=a."animalId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}

	eRows, err := s.db.Query(`SELECT e."enclosureId", e."enclosureName", e."image", e."Notes", e."habitatId"
							FROM "enclosures" e JOIN "enclosureUser" ON "enclosureUser"."enclosureId"=e."enclosureId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}

	tasks := make([]*types.Task, 0)
	for tRows.Next() {
		task := new(types.Task)
		task, err := utils.ScanRowsIntoTask(tRows)
		if err != nil {
			return err
		}

		tasks = append(tasks, task)
	}

	animals := make([]*types.Animal, 0)
	for aRows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(aRows)
		if err != nil {
			return err
		}

		animals = append(animals, animal)
	}

	enclosures := make([]*types.Enclosure, 0)
	for eRows.Next() {
		enclosure, err := utils.ScanRowsIntoEnclosures(eRows)
		if err != nil {
			return err
		}

		enclosures = append(enclosures, enclosure)
	}

	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete taskUser, taskSubject, and tasks entries
	_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, task.TaskId)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, task.TaskId)
		if err != nil {
			return err
		}
	}

	// delete animalUser and animals entries
	_, err = tx.Exec(`DELETE FROM "animalUser" WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}
	for _, animal := range animals {
		_, err = tx.Exec(`DELETE FROM "animals" WHERE "animalId" = $1`, animal.AnimalId)
		if err != nil {
			return err
		}
	}

	// delete enclosureUser and enclosures entries
	_, err = tx.Exec(`DELETE FROM "enclosureUser" WHERE "userId" = $1`, userID)
	if err != nil {
		return err
	}
	for _, enclosure := range enclosures {
		_, err = tx.Exec(`DELETE FROM "enclosures" WHERE "enclosureId" = $1`, enclosure.EnclosureId)
		if err != nil {
			return err
		}
	}

	// delete user
	_, err = tx.Exec(`DELETE FROM "users" WHERE "userId" = $1`, userID)
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

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
