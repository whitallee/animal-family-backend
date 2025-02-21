package animal

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

func (s *Store) CreateAnimal(animal types.Animal) error {
	_, err := s.db.Exec("INSERT INTO animals (animalName, speciesId, enclosureId, image, notes) VALUES (?,?,?,?,?)", animal.AnimalName, animal.SpeciesId, animal.EnclosureId, animal.Image, animal.Notes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateAnimalByUserId(animal types.Animal, userID int) error {
	// start transaction
	tx, err := s.db.Begin()

	// add animal to animals table
	tx.Exec("INSERT INTO animals (animalName, speciesId, enclosureId, image, notes) VALUES (?,?,?,?,?)", animal.AnimalName, animal.SpeciesId, animal.EnclosureId, animal.Image, animal.Notes)
	if err != nil {
		return err
	}

	// get animal id of the newly added animal
	var addedAnimalId int
	if err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&addedAnimalId); err != nil {
		return err
	}

	// add user-animal joiner to animalUser table
	tx.Exec("INSERT INTO animalUser (animalId, userID) VALUES (?,?)", addedAnimalId, userID)

	// commit transation
	tx.Commit()

	return nil
}

func (s *Store) GetAnimals() ([]*types.Animal, error) {
	rows, err := s.db.Query("SELECT * FROM animals")
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		a, err := scanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, a)
	}

	return animals, nil
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
		animal, err := scanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, animal)
	}

	return animals, nil
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
