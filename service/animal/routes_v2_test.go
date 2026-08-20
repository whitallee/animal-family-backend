package animal

import (
	"testing"
	"time"

	"github.com/whitallee/animal-family-backend/types"
)

// PUT replaces the animal's details outright: what comes back is exactly what
// was sent, with no merging against the stored record.
func TestAnimalFromUpdateCopiesPayloadVerbatim(t *testing.T) {
	dob := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	enclosureId := 4

	payload := types.UpdateAnimalV2Payload{
		AnimalName:      "Winston",
		SpeciesId:       2,
		EnclosureId:     &enclosureId,
		Image:           "https://example.test/winston.png",
		Gender:          "male",
		Dob:             dob,
		PersonalityDesc: "curious",
		DietDesc:        "omnivore",
		RoutineDesc:     "morning feed",
		ExtraNotes:      "likes boxes",
	}

	animal := animalFromUpdate(payload, 11)

	if animal.AnimalId != 11 {
		t.Errorf("id should come from the path, got %d", animal.AnimalId)
	}
	if animal.AnimalName != payload.AnimalName || animal.SpeciesId != payload.SpeciesId {
		t.Error("identifying fields were not copied")
	}
	if animal.EnclosureId == nil || *animal.EnclosureId != enclosureId {
		t.Errorf("enclosureId not copied: %v", animal.EnclosureId)
	}
	if animal.ExtraNotes != payload.ExtraNotes || animal.RoutineDesc != payload.RoutineDesc {
		t.Error("description fields were not copied")
	}
	if !animal.Dob.Equal(dob) {
		t.Errorf("dob not copied: %v", animal.Dob)
	}
}

// An omitted field is reset rather than left alone, which is what makes PUT a
// replace rather than a merge.
func TestAnimalFromUpdateResetsOmittedFields(t *testing.T) {
	animal := animalFromUpdate(types.UpdateAnimalV2Payload{
		AnimalName: "Pip",
		SpeciesId:  1,
	}, 5)

	if animal.ExtraNotes != "" || animal.PersonalityDesc != "" || animal.Gender != "" {
		t.Error("omitted fields should be zeroed, not carried over from anywhere")
	}
	if animal.EnclosureId != nil {
		t.Errorf("omitted enclosureId should be nil, got %v", *animal.EnclosureId)
	}
}

// The memorial columns are not reachable from the update payload at all. If a
// future edit adds them back, this stops compiling — which is the point: the
// edit form sends no memorial data, so a full replace that could touch those
// columns would erase a memorial message and its photos.
func TestUpdatePayloadCarriesNoMemorialState(t *testing.T) {
	animal := animalFromUpdate(types.UpdateAnimalV2Payload{AnimalName: "Pip", SpeciesId: 1}, 5)

	if animal.IsMemorialized {
		t.Error("update must not set memorial state")
	}
	if animal.LastMessage.Valid || animal.MemorialPhotos.Valid {
		t.Error("update must not set memorial content")
	}
	if !animal.MemorialDate.IsZero() {
		t.Error("update must not set the memorial date")
	}
	// UpdateAnimalDetails is what keeps those zero values from reaching the
	// database; its statement lists only the detail columns.
}
