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
	var addedAnimalId int
	err = tx.QueryRow("INSERT INTO animals (animalName, speciesId, enclosureId, image, notes) VALUES ($1, $2, $3, $4, $5) RETURNING animalId",
		animal.AnimalName, animal.SpeciesId, animal.EnclosureId, animal.Image, animal.Notes).Scan(&addedAnimalId)
	if err != nil {
		return err
	}

	// add user-animal joiner to animalUser table
	_, err = tx.Exec("INSERT INTO animalUser (animalId, userID) VALUES ($1, $2)", addedAnimalId, userID)
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
						SET animalName = $1, image = $2, notes = $3, speciesID = $4, enclosureID = $5
						WHERE animalId = $6`, animal.AnimalName, animal.Image, animal.Notes, animal.SpeciesId, animal.EnclosureId, animal.AnimalId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateAnimalOwner(oldAnimalUser types.AnimalUser, newUserId int) error {
	_, err := s.db.Exec(`UPDATE animalUser
						SET userId = $1
						WHERE animalId = $2 AND userId = $3`, newUserId, oldAnimalUser.AnimalId, oldAnimalUser.UserID)
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
							WHERE animalName = $1 AND speciesId = $2 AND userID = $3`, animalName, speciesId, userID)
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
	rows, err := s.db.Query("SELECT * FROM animalUser WHERE animalId = $1 AND userID = $2", animalId, userID)
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

func (s *Store) GetAnimalUserByAnimalId(animalId int) (*types.AnimalUser, error) {
	rows, err := s.db.Query("SELECT * FROM animalUser WHERE animalId = $1", animalId)
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
		return nil, fmt.Errorf("no owner of animal with id %d found", animalId)
	}

	return animalUser, nil
}

func (s *Store) GetAnimalById(animalId int) (*types.Animal, error) {
	rows, err := s.db.Query("SELECT * FROM animals WHERE animalId = $1", animalId)
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

func (s *Store) GetAnimalsByUserId(userID int) ([]*types.Animal, error) {
	rows, err := s.db.Query(`SELECT a.animalId, a.animalName, a.image, a.notes, a.speciesId, a.enclosureId
							FROM animals a JOIN animalUser ON animalUser.animalId=a.animalId
							WHERE userID = $1`, userID)
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
	rows, err := s.db.Query("SELECT * FROM animals WHERE enclosureID = $1", enclosureId)
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

	_, err = tx.Exec("DELETE FROM animalUser WHERE animalId = $1", animalId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM animals WHERE animalId = $1", animalId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
