package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

// WriteJSON writes v as the response body.
//
// It deliberately returns nothing. By the time it is called the status line has
// already been committed to the wire, so a caller could not act on a failure
// even if it received one — it cannot switch to a 500. Logging is the only
// meaningful response, so it happens here rather than being repeated (or, as
// errcheck otherwise pushes you toward, discarded with `_ =`) at every handler
// exit. WriteError has always worked this way; this matches it.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// 204 must not carry a body (RFC 9110 §15.3.5) and net/http refuses to
	// write one, returning "request method or response status code does not
	// allow body". That error was being produced on every successful update and
	// delete, so encoding here at all is what created the noise.
	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write %d response: %v", status, err)
	}
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

func ScanRowsIntoEnclosureUser(rows *sql.Rows) (*types.EnclosureUser, error) {
	enclosureUser := new(types.EnclosureUser)

	err := rows.Scan(
		&enclosureUser.EnclosureId,
		&enclosureUser.UserID,
	)
	if err != nil {
		return nil, err
	}

	return enclosureUser, nil
}

func ScanRowsIntoAnimals(rows *sql.Rows) (*types.Animal, error) {
	animal := new(types.Animal)

	err := rows.Scan(
		&animal.AnimalId,
		&animal.AnimalName,
		&animal.Image,
		&animal.ExtraNotes,
		&animal.SpeciesId,
		&animal.EnclosureId,
		&animal.Gender,
		&animal.Dob,
		&animal.PersonalityDesc,
		&animal.DietDesc,
		&animal.RoutineDesc,
		&animal.IsMemorialized,
		&animal.LastMessage,
		&animal.MemorialPhotos,
		&animal.MemorialDate,
	)
	if err != nil {
		return nil, err
	}

	return animal, nil
}

func ScanRowsIntoAnimalUser(rows *sql.Rows) (*types.AnimalUser, error) {
	animalUser := new(types.AnimalUser)

	err := rows.Scan(
		&animalUser.AnimalId,
		&animalUser.UserID,
	)
	if err != nil {
		return nil, err
	}

	return animalUser, nil
}

func ScanRowsIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)

	err := rows.Scan(
		&task.TaskId,
		&task.TaskName,
		&task.TaskDesc,
		&task.Complete,
		&task.LastCompleted,
		&task.RepeatIntervHours,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func ScanRowsIntoTaskWithSubject(rows *sql.Rows) (*types.TaskWithSubject, error) {
	task := new(types.TaskWithSubject)

	err := rows.Scan(
		&task.TaskId,
		&task.TaskName,
		&task.TaskDesc,
		&task.Complete,
		&task.LastCompleted,
		&task.RepeatIntervHours,
		&task.AnimalId,
		&task.EnclosureId,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func ScanRowsIntoTaskUser(rows *sql.Rows) (*types.TaskUser, error) {
	taskUser := new(types.TaskUser)

	err := rows.Scan(
		&taskUser.TaskId,
		&taskUser.UserID,
	)
	if err != nil {
		return nil, err
	}

	return taskUser, nil
}
