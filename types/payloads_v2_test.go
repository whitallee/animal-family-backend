package types

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// Anything that can be created must stay editable.
//
// v1 required gender, dob and the description fields on update while leaving
// them optional on create, so an animal saved without them could never be
// updated again. Nothing caught it because the two payload types were written
// separately and never compared.
func TestUpdateAnimalAcceptsWhatCreateAccepts(t *testing.T) {
	validate := validator.New()

	// The minimum a create will accept: everything else left at its zero value.
	create := CreateAnimalV2Payload{AnimalName: "Pip", SpeciesId: 1}
	if err := validate.Struct(create); err != nil {
		t.Fatalf("this payload is supposed to be creatable: %v", err)
	}

	update := UpdateAnimalV2Payload{AnimalName: create.AnimalName, SpeciesId: create.SpeciesId}
	if err := validate.Struct(update); err != nil {
		t.Errorf("an animal created with these fields could not be updated: %v", err)
	}
}

func TestUpdateAnimalStillRequiresIdentifyingFields(t *testing.T) {
	validate := validator.New()

	cases := map[string]UpdateAnimalV2Payload{
		"missing name":    {SpeciesId: 1},
		"missing species": {AnimalName: "Pip"},
		"species id zero": {AnimalName: "Pip", SpeciesId: 0},
	}

	for name, payload := range cases {
		t.Run(name, func(t *testing.T) {
			if err := validate.Struct(payload); err == nil {
				t.Error("expected validation to fail")
			}
		})
	}
}

// The same symmetry for enclosures, which already held but is worth pinning so
// a later edit cannot reintroduce the asymmetry.
func TestUpdateEnclosureAcceptsWhatCreateAccepts(t *testing.T) {
	validate := validator.New()

	create := CreateEnclosureV2Payload{EnclosureName: "Terrarium", HabitatId: 1}
	if err := validate.Struct(create); err != nil {
		t.Fatalf("this payload is supposed to be creatable: %v", err)
	}

	update := UpdateEnclosureV2Payload{EnclosureName: create.EnclosureName, HabitatId: create.HabitatId}
	if err := validate.Struct(update); err != nil {
		t.Errorf("an enclosure created with these fields could not be updated: %v", err)
	}
}
