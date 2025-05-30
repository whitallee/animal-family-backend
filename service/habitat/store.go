package habitat

import (
	"database/sql"
	"fmt"

	"github.com/whitallee/animal-family-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateHabitat(habitat types.Habitat) error {
	_, err := s.db.Exec(`INSERT INTO "habitats" ("habitatId", "habitatName", "habitatDesc", "image", "humidity", "dayTempRange", "nightTempRange") VALUES ($1,$2,$3,$4,$5,$6,$7)`, habitat.HabitatId, habitat.HabitatName, habitat.HabitatDesc, habitat.Image, habitat.Humidity, habitat.DayTempRange, habitat.NightTempRange)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateHabitat(habitat types.Habitat) error {
	_, err := s.db.Exec(`UPDATE "habitats"
						SET "habitatName" = $1, "habitatDesc" = $2, "image" = $3, "humidity" = $4, "dayTempRange" = $5, "nightTempRange" = $6
						WHERE "habitatId" = $7`, habitat.HabitatName, habitat.HabitatDesc, habitat.Image, habitat.Humidity, habitat.DayTempRange, habitat.NightTempRange, habitat.HabitatId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetHabitats() ([]*types.Habitat, error) {
	rows, err := s.db.Query(`SELECT * FROM "habitats"`)
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

func (s *Store) GetHabitatByName(habName string) (*types.Habitat, error) {
	rows, err := s.db.Query(`SELECT * FROM "habitats" WHERE "habitatName" = $1`, habName)
	if err != nil {
		return nil, err
	}

	habitat := new(types.Habitat)
	for rows.Next() {
		habitat, err = scanRowsIntoHabitats(rows)
		if err != nil {
			return nil, err
		}
	}

	if habitat.HabitatId == 0 {
		return nil, fmt.Errorf("habitat not found")
	}

	return habitat, nil
}

func (s *Store) GetHabitatById(habId int) (*types.Habitat, error) { // not used in any handler functions yet
	rows, err := s.db.Query(`SELECT * FROM "habitats" WHERE "habitatId" = $1`, habId)
	if err != nil {
		return nil, err
	}

	habitat := new(types.Habitat)
	for rows.Next() {
		habitat, err = scanRowsIntoHabitats(rows)
		if err != nil {
			return nil, err
		}
	}

	if habitat.HabitatId == 0 {
		return nil, fmt.Errorf("habitat not found")
	}

	return habitat, nil
}

func (s *Store) DeleteHabitatById(habitatId int) error {
	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// update enclosures habitat to "No Habitat"
	_, err = tx.Exec(`UPDATE "enclosures" SET "habitatId" = 0 WHERE "habitatId" = $1`, habitatId)
	if err != nil {
		return err
	}

	// update species habitat to "No Habitat"
	_, err = tx.Exec(`UPDATE "species" SET "habitatId" = 0 WHERE "habitatId" = $1`, habitatId)
	if err != nil {
		return err
	}

	// delete habitat
	_, err = tx.Exec(`DELETE FROM "habitats" WHERE "habitatId" = $1`, habitatId)
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
