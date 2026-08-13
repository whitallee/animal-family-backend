package types

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
