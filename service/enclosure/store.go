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
