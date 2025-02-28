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
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?,?,?,?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
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
	rows, err := s.db.Query("SELECT * FROM users WHERE userId = ?", id)
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
	// get animals and enclosures from userID
	aRows, err := s.db.Query(`SELECT a.animalId, a.animalName, a.image, a.notes, a.speciesId, a.enclosureId
							FROM animals a JOIN animalUser ON animalUser.animalId=a.animalId
							WHERE userID = ?`, userID)
	if err != nil {
		return err
	}

	eRows, err := s.db.Query(`SELECT e.enclosureId, e.enclosureName, e.image, e.Notes, e.habitatId
							FROM enclosures e JOIN enclosureUser ON enclosureUser.enclosureId=e.enclosureId
							WHERE userID = ?`, userID)
	if err != nil {
		return err
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

	// delete animalUser and animals entries
	_, err = tx.Exec("DELETE FROM animalUser WHERE userID = ?", userID)
	if err != nil {
		return err
	}
	for _, animal := range animals {
		_, err = tx.Exec("DELETE FROM animals WHERE animalId = ?", animal.AnimalId)
		if err != nil {
			return err
		}
	}

	// delete enclosureUser and enclosures entries
	_, err = tx.Exec("DELETE FROM enclosureUser WHERE userID = ?", userID)
	if err != nil {
		return err
	}
	for _, enclosure := range enclosures {
		_, err = tx.Exec("DELETE FROM enclosures WHERE enclosureId = ?", enclosure.EnclosureId)
		if err != nil {
			return err
		}
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
