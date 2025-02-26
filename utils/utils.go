package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/whitallee/animal-family-backend/types"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ScanRowsIntoEnclosures(rows *sql.Rows) (*types.Enclosure, error) {
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

func ScanRowsIntoAnimals(rows *sql.Rows) (*types.Animal, error) {
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
