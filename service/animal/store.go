package animal

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

func (s *Store) CreateAnimal(animal types.Animal, userID int) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// add animal to animals table
	_, err = tx.Exec("INSERT INTO animals (animalName, speciesId, enclosureId, image, notes) VALUES (?,?,?,?,?)", animal.AnimalName, animal.SpeciesId, animal.EnclosureId, animal.Image, animal.Notes)
	if err != nil {
		return err
	}

	// get animal id of the newly added animal
	var addedAnimalId int
	if err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&addedAnimalId); err != nil {
		return err
	}

	// add user-animal joiner to animalUser table
	_, err = tx.Exec("INSERT INTO animalUser (animalId, userID) VALUES (?,?)", addedAnimalId, userID)
	if err != nil {
		return err
	}

	// commit transation
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateAnimal(animal types.Animal) error {
	_, err := s.db.Exec(`UPDATE animals
						SET animalName = ?, image = ?, notes = ?, speciesID = ?, enclosureID = ?
						WHERE animalId = ?`, animal.AnimalName, animal.Image, animal.Notes, animal.SpeciesId, animal.EnclosureId, animal.AnimalId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAnimals() ([]*types.Animal, error) {
	rows, err := s.db.Query("SELECT * FROM animals")
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		a, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, a)
	}

	return animals, nil
}

func (s *Store) GetAnimalByNameAndSpeciesWithUserId(animalName string, speciesId int, userID int) (*types.Animal, error) {
	rows, err := s.db.Query(`SELECT a.animalId, a.animalName, a.image, a.notes, a.speciesId, a.enclosureId
							FROM animals a JOIN animalUser ON animalUser.animalId=a.animalId
							WHERE animalName = ? AND speciesId = ? AND userID = ?`, animalName, speciesId, userID)
	if err != nil {
		return nil, err
	}

	animal := new(types.Animal)
	for rows.Next() {
		animal, err = utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}
	}

	if animal.AnimalId == 0 {
		return nil, fmt.Errorf("animal not found")
	}

	return animal, nil
}

func (s *Store) GetAnimalUserByIds(animalId int, userID int) (*types.AnimalUser, error) {
	rows, err := s.db.Query("SELECT * FROM animalUser WHERE animalId = ? AND userID = ?", animalId, userID)
	if err != nil {
		return nil, err
	}

	animalUser := new(types.AnimalUser)
	for rows.Next() {
		animalUser, err = utils.ScanRowsIntoAnimalUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if animalUser.AnimalId == 0 && animalUser.UserID == 0 {
		return nil, fmt.Errorf("no ownership found between user and animal")
	}

	return animalUser, nil
}

func (s *Store) GetAnimalsByUserId(userID int) ([]*types.Animal, error) {
	rows, err := s.db.Query(`SELECT a.animalId, a.animalName, a.image, a.notes, a.speciesId, a.enclosureId
							FROM animals a JOIN animalUser ON animalUser.animalId=a.animalId
							WHERE userID = ?`, userID)
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, animal)
	}

	return animals, nil
}

func (s *Store) GetAnimalsByEnclosureId(enclosureId int) ([]*types.Animal, error) {
	rows, err := s.db.Query("SELECT * FROM animals WHERE enclosureID = ?", enclosureId)
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, animal)
	}

	return animals, nil
}

func (s *Store) DeleteAnimalById(animalId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM animalUser WHERE animalId = ?", animalId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM animals WHERE animalId = ?", animalId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
