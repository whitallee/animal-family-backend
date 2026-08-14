package types

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// AnimalResponse exists to describe the wire format that Animal.MarshalJSON
// already produces, so the two must encode identically. If they diverge, the
// generated client's types stop matching what the server actually sends —
// silently, since nothing would fail to compile.
func TestAnimalResponseMatchesAnimalMarshalJSON(t *testing.T) {
	dob := time.Date(2020, 3, 14, 0, 0, 0, 0, time.UTC)
	memorialDate := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	cases := map[string]Animal{
		"populated nullables": {
			AnimalId:        7,
			AnimalName:      "Winston",
			SpeciesId:       3,
			EnclosureId:     new(2),
			Image:           "https://example.test/winston.png",
			Gender:          "male",
			Dob:             dob,
			PersonalityDesc: "curious",
			DietDesc:        "omnivore",
			RoutineDesc:     "morning feed",
			ExtraNotes:      "likes boxes",
			IsMemorialized:  true,
			LastMessage:     sql.NullString{String: "we miss you", Valid: true},
			MemorialPhotos:  sql.NullString{String: `["a.png","b.png"]`, Valid: true},
			MemorialDate:    memorialDate,
		},
		"null nullables": {
			AnimalId:       8,
			AnimalName:     "Pip",
			SpeciesId:      1,
			EnclosureId:    nil,
			Image:          "https://example.test/pip.png",
			Gender:         "female",
			Dob:            dob,
			IsMemorialized: false,
			LastMessage:    sql.NullString{Valid: false},
			MemorialPhotos: sql.NullString{Valid: false},
			MemorialDate:   memorialDate,
		},
		"empty photo array": {
			AnimalId:       9,
			AnimalName:     "Moss",
			SpeciesId:      2,
			EnclosureId:    new(5),
			MemorialPhotos: sql.NullString{String: `[]`, Valid: true},
			Dob:            dob,
			MemorialDate:   memorialDate,
		},
	}

	for name, animal := range cases {
		t.Run(name, func(t *testing.T) {
			legacyJSON, err := json.Marshal(animal)
			if err != nil {
				t.Fatalf("marshal Animal: %v", err)
			}

			response, err := NewAnimalResponse(&animal)
			if err != nil {
				t.Fatalf("NewAnimalResponse: %v", err)
			}

			responseJSON, err := json.Marshal(response)
			if err != nil {
				t.Fatalf("marshal AnimalResponse: %v", err)
			}

			// Compared as decoded maps because JSON object key order is not
			// significant and the two types declare fields in a different order.
			var legacyFields, responseFields map[string]any
			if err := json.Unmarshal(legacyJSON, &legacyFields); err != nil {
				t.Fatalf("decode Animal JSON: %v", err)
			}
			if err := json.Unmarshal(responseJSON, &responseFields); err != nil {
				t.Fatalf("decode AnimalResponse JSON: %v", err)
			}

			if !reflect.DeepEqual(legacyFields, responseFields) {
				t.Errorf("wire format changed\n Animal: %s\nResponse: %s", legacyJSON, responseJSON)
			}
		})
	}
}

// The two fields the spec generator cannot see on Animal, because they are
// tagged `json:"-"` and only appear via MarshalJSON.
func TestAnimalResponseCarriesMemorialFields(t *testing.T) {
	animal := Animal{
		AnimalId:       1,
		LastMessage:    sql.NullString{String: "goodbye", Valid: true},
		MemorialPhotos: sql.NullString{String: `["x.png"]`, Valid: true},
	}

	response, err := NewAnimalResponse(&animal)
	if err != nil {
		t.Fatalf("NewAnimalResponse: %v", err)
	}

	if response.LastMessage == nil || *response.LastMessage != "goodbye" {
		t.Errorf("lastMessage not carried through: %v", response.LastMessage)
	}
	if !reflect.DeepEqual(response.MemorialPhotos, []string{"x.png"}) {
		t.Errorf("memorialPhotos not carried through: %v", response.MemorialPhotos)
	}
}

func TestNewAnimalResponseRejectsMalformedPhotos(t *testing.T) {
	animal := Animal{
		AnimalId:       4,
		MemorialPhotos: sql.NullString{String: "not json", Valid: true},
	}

	if _, err := NewAnimalResponse(&animal); err == nil {
		t.Error("expected an error for a non-JSON memorialPhotos column")
	}
}

// A nil slice encodes as `null`, but list endpoints are typed as returning an
// array, so an empty result must still encode as [].
func TestNewAnimalResponsesEncodesEmptyListAsArray(t *testing.T) {
	responses, err := NewAnimalResponses(nil)
	if err != nil {
		t.Fatalf("NewAnimalResponses: %v", err)
	}

	encoded, err := json.Marshal(responses)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(encoded) != "[]" {
		t.Errorf("expected [], got %s", encoded)
	}
}
