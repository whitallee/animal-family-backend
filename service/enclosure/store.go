package enclosure

import (
	"database/sql"

	"github.com/whitallee/animal-family-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateEnclosure(enclosure types.Enclosure) error {
	_, err := s.db.Exec("INSERT INTO enclosures (enclosureName, image, notes, habitatId) VALUES (?,?,?,?)", enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateEnclosureByUserId(enclosure types.Enclosure, userID int) error {
	tx, err := s.db.Begin()

	tx.Exec("INSERT INTO enclosures (enclosureName, image, notes, habitatId) VALUES (?,?,?,?)", enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId)
	if err != nil {
		return err
	}

	var addedEnclosureId int
	if err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&addedEnclosureId); err != nil {
		return err
	}

	tx.Exec("INSERT INTO enclosureUser (enclosureId, userID) VALUES (?,?)", addedEnclosureId, userID)

	tx.Commit()

	return nil
}

func (s *Store) CreateEnclosureWithAnimalsByUserId(enclosure types.Enclosure, animalIds []int, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO enclosures (enclosureName, image, notes, habitatId) VALUES (?,?,?,?)", enclosure.EnclosureName, enclosure.Image, enclosure.Notes, enclosure.HabitatId)
	if err != nil {
		return err
	}

	var addedEnclosureId int
	if err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&addedEnclosureId); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO enclosureUser (enclosureId, userID) VALUES (?,?)", addedEnclosureId, userID); err != nil {
		return err
	}

	for _, animalId := range animalIds {
		if _, err := tx.Exec("UPDATE animals SET enclosureID = ? WHERE animalId = ?", addedEnclosureId, animalId); err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}

func (s *Store) GetEnclosures() ([]*types.Enclosure, error) {
	rows, err := s.db.Query("SELECT * FROM enclosures")
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		s, err := scanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}

		enclosures = append(enclosures, s)
	}

	return enclosures, nil
}

func (s *Store) GetEnclosuresByUserId(userID int) ([]*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT e.enclosureId, e.enclosureName, e.image, e.Notes, e.habitatId
							FROM enclosures e JOIN enclosureUser ON enclosureUser.enclosureId=e.enclosureId
							WHERE userID = ?`, userID)
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		enclosure, err := scanRowsIntoEnclosures(rows)
		if err != nil {
			return nil, err
		}

		enclosures = append(enclosures, enclosure)
	}

	return enclosures, nil
}

func (s *Store) GetEnclosureByIdWithUserId(enclosureId int, userID int) (*types.Enclosure, error) {
	rows, err := s.db.Query(`SELECT e.enclosureId, e.enclosureName, e.image, e.Notes, e.habitatId
							FROM enclosures e JOIN enclosureUser ON enclosureUser.enclosureId=e.enclosureId
							WHERE userID = ? AND e.enclosureId = ?`, userID, enclosureId)
	if err != nil {
		return nil, err
	}

	enclosures := make([]*types.Enclosure, 0)
	for rows.Next() {
		enclosure, err := scanRowsIntoEnclosures(rows)
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

func (s *Store) DeleteEnclosureAndAnimalsByIdWithUserId(enclosureId int, userID int) error {
	// get animals fom enclosures
	rows, err := s.db.Query(`SELECT a.animalId, a.animalName, a.image, a.notes, a.speciesId, a.enclosureId
							FROM animals a JOIN animalUser ON animalUser.animalId=a.animalId
							WHERE userID = ? AND enclosureID = ?`, userID, enclosureId)
	if err != nil {
		return err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := scanRowsIntoAnimals(rows)
		if err != nil {
			return err
		}

		animals = append(animals, animal)
	}

	// start deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete from animalUser and animals
	for _, animal := range animals {
		_, err = tx.Exec("DELETE FROM animalUser WHERE animalId = ? AND userID = ?", animal.AnimalId, userID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM animals WHERE animalId = ?", animal.AnimalId)
		if err != nil {
			return err
		}
	}

	// delete from enclosureUser and enclosures
	_, err = tx.Exec("DELETE FROM enclosureUser WHERE enclosureId = ? AND userId = ?", enclosureId, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM enclosures WHERE enclosureId = ?", enclosureId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoEnclosures(rows *sql.Rows) (*types.Enclosure, error) {
	enclosures := new(types.Enclosure)

	err := rows.Scan(
		&enclosures.EnclosureId,
		&enclosures.EnclosureName,
		&enclosures.Image,
		&enclosures.Notes,
		&enclosures.HabitatId,
	)
	if err != nil {
		return nil, err
	}

	return enclosures, nil
}

func scanRowsIntoAnimals(rows *sql.Rows) (*types.Animal, error) {
	animal := new(types.Animal)

	err := rows.Scan(
		&animal.AnimalId,
		&animal.AnimalName,
		&animal.Image,
		&animal.Notes,
		&animal.SpeciesId,
		&animal.EnclosureId,
	)
	if err != nil {
		return nil, err
	}

	return animal, nil
}
