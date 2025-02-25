package habitat

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

func (s *Store) CreateHabitat(habitat types.Habitat) error {
	_, err := s.db.Exec("INSERT INTO habitats (habitatId, habitatName, habitatDesc, image, humidity, dayTempRange, nightTempRange) VALUES (?,?,?,?,?,?,?)", habitat.HabitatId, habitat.HabitatName, habitat.HabitatDesc, habitat.Image, habitat.Humidity, habitat.DayTempRange, habitat.NightTempRange)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetHabitats() ([]*types.Habitat, error) {
	rows, err := s.db.Query("SELECT * FROM habitats")
	if err != nil {
		return nil, err
	}

	habitats := make([]*types.Habitat, 0)
	for rows.Next() {
		h, err := scanRowsIntoHabitats(rows)
		if err != nil {
			return nil, err
		}

		habitats = append(habitats, h)
	}

	return habitats, nil
}

func (s *Store) DeleteHabitatById(habitatId int) error {
	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// update enclsoures habitat to "No Habitat"
	_, err = tx.Exec("UPDATE enclosures SET habitatId = 0 WHERE habitatId = ?", habitatId)
	if err != nil {
		return err
	}

	// update species habitat to "No Habitat"
	_, err = tx.Exec("UPDATE species SET habitatId = 0 WHERE habitatId = ?", habitatId)
	if err != nil {
		return err
	}

	// delete habitat
	_, err = tx.Exec("DELETE FROM habitats WHERE habitatId = ?", habitatId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoHabitats(rows *sql.Rows) (*types.Habitat, error) {
	habitat := new(types.Habitat)

	err := rows.Scan(
		&habitat.HabitatId,
		&habitat.HabitatName,
		&habitat.HabitatDesc,
		&habitat.Image,
		&habitat.Humidity,
		&habitat.DayTempRange,
		&habitat.NightTempRange,
	)
	if err != nil {
		return nil, err
	}

	return habitat, nil
}
