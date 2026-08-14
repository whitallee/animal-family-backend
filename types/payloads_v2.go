package types

import "time"

// Request payloads for the v2 API.
//
// v2 takes resource IDs from the URL path, so these mirror their v1
// counterparts minus the ID field. The "V2" in these names disambiguates them
// from the v1 types while both exist; cmd/specgen strips the Go package prefix
// from schema names, and the suffix can be dropped here once v1 is retired.

// UpdateHabitatV2Payload is the body of PUT /habitats/{id}.
type UpdateHabitatV2Payload struct {
	HabitatName    string `json:"habitatName" validate:"required"`
	HabitatDesc    string `json:"habitatDesc" validate:"required"`
	Image          string `json:"image" validate:"required"`
	Humidity       string `json:"humidity" validate:"required"`
	DayTempRange   string `json:"dayTempRange" validate:"required"`
	NightTempRange string `json:"nightTempRange" validate:"required"`
}

// CreateEnclosureV2Payload is the body of POST /enclosures.
//
// AnimalIds folds v1's separate /enclosure/withanimals route into this one:
// omit it to create an empty enclosure, supply it to populate the enclosure at
// the same time.
type CreateEnclosureV2Payload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=1"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
	AnimalIds     []int  `json:"animalIds" validate:"omitempty,dive,min=1"`
}

// UpdateEnclosureV2Payload is the body of PUT /enclosures/{id}.
type UpdateEnclosureV2Payload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=1"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
}

// CreateAnimalV2Payload is the body of POST /animals.
type CreateAnimalV2Payload struct {
	AnimalName      string    `json:"animalName" validate:"required"`
	SpeciesId       int       `json:"speciesId" validate:"required,min=1"`
	EnclosureId     *int      `json:"enclosureId" validate:"omitempty,min=1" extensions:"x-nullable"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender"`
	Dob             time.Time `json:"dob"`
	PersonalityDesc string    `json:"personalityDesc"`
	DietDesc        string    `json:"dietDesc"`
	RoutineDesc     string    `json:"routineDesc"`
	ExtraNotes      string    `json:"extraNotes"`
}

// UpdateAnimalV2Payload is the body of PUT /animals/{id}.
//
// The memorial fields are pointers so that omitting them leaves the stored
// values untouched, matching v1's UpdateAnimalPayload behaviour.
type UpdateAnimalV2Payload struct {
	AnimalName      string    `json:"animalName" validate:"required"`
	SpeciesId       int       `json:"speciesId" validate:"required,min=1"`
	EnclosureId     *int      `json:"enclosureId" validate:"omitempty,min=1" extensions:"x-nullable"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender" validate:"required"`
	Dob             time.Time `json:"dob" validate:"required"`
	PersonalityDesc string    `json:"personalityDesc" validate:"required"`
	DietDesc        string    `json:"dietDesc" validate:"required"`
	RoutineDesc     string    `json:"routineDesc" validate:"required"`
	ExtraNotes      string    `json:"extraNotes" validate:"required"`
	IsMemorialized  *bool     `json:"isMemorialized" extensions:"x-nullable"`
	LastMessage     *string   `json:"lastMessage" extensions:"x-nullable"`
	MemorialPhotos  []string  `json:"memorialPhotos" extensions:"x-nullable"`
	MemorialDate    time.Time `json:"memorialDate"`
}

// UpdateSpeciesV2Payload is the body of PUT /species/{id}.
type UpdateSpeciesV2Payload struct {
	ComName            string `json:"comName" validate:"required"`
	SciName            string `json:"sciName" validate:"required"`
	SpeciesDesc        string `json:"speciesDesc" validate:"required"`
	Image              string `json:"image" validate:"required"`
	HabitatId          int    `json:"habitatId" validate:"required,min=1"`
	BaskTemp           string `json:"baskTemp" validate:"required"`
	Diet               string `json:"diet" validate:"required"`
	Sociality          string `json:"sociality" validate:"required"`
	Lifespan           string `json:"lifespan" validate:"required"`
	Size               string `json:"size" validate:"required"`
	Weight             string `json:"weight" validate:"required"`
	ConservationStatus string `json:"conservationStatus" validate:"required"`
	ExtraCare          string `json:"extraCare" validate:"required"`
}
