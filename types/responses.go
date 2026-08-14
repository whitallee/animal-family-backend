package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// Response types used by the v2 API.
//
// These exist so the OpenAPI generator has concrete types to reference. v1
// wrote several responses as ad hoc map[string]interface{} literals, which
// produce correct JSON but leave nothing for the spec (and therefore the
// generated TypeScript client) to describe.

// ErrorResponse is the body written by utils.WriteError for every non-2xx
// response. Its shape must stay in sync with that function.
type ErrorResponse struct {
	Error string `json:"error" example:"admin access required"`
}

// AnimalResponse is the wire representation of an animal.
//
// Animal itself cannot be used here. Its LastMessage and MemorialPhotos fields
// are sql.NullString tagged `json:"-"`, and its custom MarshalJSON synthesises
// `lastMessage` and `memorialPhotos` at encode time. The spec generator reflects
// struct tags and cannot see marshalling code, so generating from Animal would
// silently omit both fields — exactly the two the memorial feature depends on.
//
// Declaring the wire shape explicitly also keeps database types out of the API
// contract, and gives nullable fields somewhere to carry `x-nullable` (Swagger
// 2.0 has no nullable keyword; specgen's converter turns the extension into
// OpenAPI 3's `nullable: true`).
//
// Build one with NewAnimalResponse rather than by hand.
type AnimalResponse struct {
	AnimalId        int       `json:"animalId"`
	AnimalName      string    `json:"animalName"`
	SpeciesId       int       `json:"speciesId"`
	EnclosureId     *int      `json:"enclosureId" extensions:"x-nullable"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender"`
	Dob             time.Time `json:"dob"`
	PersonalityDesc string    `json:"personalityDesc"`
	DietDesc        string    `json:"dietDesc"`
	RoutineDesc     string    `json:"routineDesc"`
	ExtraNotes      string    `json:"extraNotes"`
	IsMemorialized  bool      `json:"isMemorialized"`
	LastMessage     *string   `json:"lastMessage" extensions:"x-nullable"`
	MemorialPhotos  []string  `json:"memorialPhotos" extensions:"x-nullable"`
	MemorialDate    time.Time `json:"memorialDate"`
}

// NewAnimalResponse converts a stored animal into its wire representation.
//
// It returns an error only when MemorialPhotos holds text that is not a JSON
// array, which would mean the column was written by something other than the
// memorialise flow.
func NewAnimalResponse(a *Animal) (AnimalResponse, error) {
	response := AnimalResponse{
		AnimalId:        a.AnimalId,
		AnimalName:      a.AnimalName,
		SpeciesId:       a.SpeciesId,
		EnclosureId:     a.EnclosureId,
		Image:           a.Image,
		Gender:          a.Gender,
		Dob:             a.Dob,
		PersonalityDesc: a.PersonalityDesc,
		DietDesc:        a.DietDesc,
		RoutineDesc:     a.RoutineDesc,
		ExtraNotes:      a.ExtraNotes,
		IsMemorialized:  a.IsMemorialized,
		MemorialDate:    a.MemorialDate,
	}

	if a.LastMessage.Valid {
		message := a.LastMessage.String
		response.LastMessage = &message
	}

	if a.MemorialPhotos.Valid {
		var photos []string
		if err := json.Unmarshal([]byte(a.MemorialPhotos.String), &photos); err != nil {
			return AnimalResponse{}, fmt.Errorf("animal %d has malformed memorialPhotos: %w", a.AnimalId, err)
		}
		response.MemorialPhotos = photos
	}

	return response, nil
}

// NewAnimalResponses converts a list of stored animals for a collection route.
func NewAnimalResponses(animals []*Animal) ([]AnimalResponse, error) {
	// Never nil: a nil slice encodes as `null`, and the generated client types
	// list endpoints as returning an array.
	responses := make([]AnimalResponse, 0, len(animals))

	for _, a := range animals {
		response, err := NewAnimalResponse(a)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}
