package types

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// User.Phone is a sql.NullString with a live `json:"phone"` tag, so encoding a
// User directly puts {"String":...,"Valid":...} on the wire where the frontend
// declared a plain string. UserResponse is what makes the contract honest.
func TestUserResponseEncodesPhoneAsNullableString(t *testing.T) {
	created := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)

	t.Run("phone present", func(t *testing.T) {
		user := User{
			ID: 1, FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.test",
			Phone: sql.NullString{String: "555-0100", Valid: true}, CreatedAt: created,
		}

		encoded, err := json.Marshal(NewUserResponse(&user))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		var decoded map[string]any
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatalf("decode: %v", err)
		}

		if decoded["phone"] != "555-0100" {
			t.Errorf("expected phone to be a plain string, got %#v", decoded["phone"])
		}
	})

	t.Run("phone absent encodes as null", func(t *testing.T) {
		user := User{ID: 2, Email: "b@example.test", Phone: sql.NullString{Valid: false}, CreatedAt: created}

		encoded, err := json.Marshal(NewUserResponse(&user))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		var decoded map[string]any
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatalf("decode: %v", err)
		}

		value, present := decoded["phone"]
		if !present {
			t.Error("phone key should still be present, as the schema marks it required")
		}
		if value != nil {
			t.Errorf("expected null, got %#v", value)
		}
	})
}

// A password must never appear in a response. UserResponse has no such field,
// so this guards against one being added later.
func TestUserResponseNeverIncludesPassword(t *testing.T) {
	user := User{ID: 3, Email: "c@example.test", Password: "super-secret-hash"}

	encoded, err := json.Marshal(NewUserResponse(&user))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if strings.Contains(string(encoded), "super-secret-hash") || strings.Contains(string(encoded), "password") {
		t.Errorf("password material leaked into the response: %s", encoded)
	}
}

// The stored subscription carries the p256dh and auth encryption keys; the
// response type deliberately omits them.
func TestPushSubscriptionResponseOmitsEncryptionKeys(t *testing.T) {
	subscription := PushSubscription{
		SubscriptionId: 1,
		UserID:         4,
		Endpoint:       "https://push.example/abc",
		P256dh:         "p256dh-secret-material",
		Auth:           "auth-secret-material",
		UserAgent:      "test-agent",
	}

	encoded, err := json.Marshal(NewPushSubscriptionResponse(&subscription))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	for _, secret := range []string{"p256dh-secret-material", "auth-secret-material"} {
		if strings.Contains(string(encoded), secret) {
			t.Errorf("encryption key material leaked: %s", encoded)
		}
	}
}

func TestNewPushSubscriptionResponsesEncodesEmptyListAsArray(t *testing.T) {
	encoded, err := json.Marshal(NewPushSubscriptionResponses(nil))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(encoded) != "[]" {
		t.Errorf("expected [], got %s", encoded)
	}
}

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
