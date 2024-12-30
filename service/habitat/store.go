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

func (s *Store) CreateSpecies(species types.Species) error {
	_, err := s.db.Exec("INSERT INTO species (comName, sciName, speciesDesc, image, habitatId, baskTemp, diet, sociality, extraCare) VALUES (?,?,?,?,?,?,?,?,?)", species.ComName, species.SciName, species.SpeciesDesc, species.Image, species.HabitatId, species.BaskTemp, species.Diet, species.Sociality, species.ExtraCare)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSpecies() ([]*types.Species, error) {
	rows, err := s.db.Query("SELECT * FROM species")
	if err != nil {
		return nil, err
	}

	species := make([]*types.Species, 0)
	for rows.Next() {
		s, err := scanRowsIntoSpecies(rows)
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
		&species.SpeciesID,
		&species.ComName,
		&species.SciName,
		&species.SpeciesDesc,
		&species.Image,
		&species.HabitatId,
		&species.BaskTemp,
		&species.Diet,
		&species.Sociality,
		&species.ExtraCare,
	)
	if err != nil {
		return nil, err
	}

	return species, nil
}
