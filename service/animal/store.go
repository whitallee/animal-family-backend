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

func scanRowsIntoAnimals(rows *sql.Rows) (*types.Animal, error) {
	animal := new(types.Animal)

	err := rows.Scan(
		&animal.AnimalId,
		&animal.AnimalName,
		&animal.SpeciesId,
		&animal.EnclosureId,
		&animal.Image,
		&animal.Notes,
	)
	if err != nil {
		return nil, err
	}

	return animal, nil
}
