package species

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

func (s *Store) GetSpecies() ([]*types.Species, error) {
	rows, err := s.db.Query("SELECT * FROM species")
	if err != nil {
		return nil, err
	}

	species := make([]*types.Species, 0)
	for rows.Next() {
		s, err := scanRowsIntoSpecies(rows) //CHANGE to SPECIES
		if err != nil {
			return nil, err
		}

		species = append(species, s)
	}

	return species, nil
}

func scanRowsIntoSpecies(rows *sql.Rows) (*types.Species, error) {
	species := new(types.Species)

	err := rows.Scan(
		&species.ID,
		&species.Name,
		&species.SciName,
		&species.Description,
		&species.Image,
		&species.Habitat,
		&species.Diet,
		&species.Sociality,
	)
	if err != nil {
		return nil, err
	}

	return species, nil
}
