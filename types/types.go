package types

import "time"

// User-related types
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Species-related types
type SpeciesStore interface {
	GetSpecies() ([]*Species, error)
	CreateSpecies(Species) error
}

type Species struct {
	SpeciesID   int    `json:"speciesId"`
	ComName     string `json:"comName"`
	SciName     string `json:"sciName"`
	SpeciesDesc string `json:"speciesDesc"`
	Image       string `json:"image"`
	HabitatId   int    `json:"habitatId"`
	BaskTemp    string `json:"baskTemp"`
	Diet        string `json:"diet"`
	Sociality   string `json:"sociality"`
	ExtraCare   string `json:"extraCare"`
}

type CreateSpeciesPayload struct {
	ComName     string `json:"comName"`
	SciName     string `json:"sciName"`
	SpeciesDesc string `json:"speciesDesc"`
	Image       string `json:"image"`
	HabitatId   int    `json:"habitatId"`
	BaskTemp    string `json:"baskTemp"`
	Diet        string `json:"diet"`
	Sociality   string `json:"sociality"`
	ExtraCare   string `json:"extraCare"`
}

// Habitat-related types
type Habitat struct { //HABITAT SERVICE NEEDED
	HabitatId      int    `json:"habitatId"`
	HabitatName    string `json:"habitatName"`
	HabitatDesc    string `json:"habitatDesc"`
	Image          string `json:"image"`
	Humidity       string `json:"humidity"`
	DayTempRange   string `json:"dayTempRange"`
	NightTempRange string `json:"nightTempRange"`
}

type PostHabitatPayload struct { //HABITAT SERVICE NEEDED
	HabitatId      int    `json:"habitatId"`
	HabitatName    string `json:"habitatName"`
	HabitatDesc    string `json:"habitatDesc"`
	Image          string `json:"image"`
	Humidity       string `json:"humidity"`
	DayTempRange   string `json:"dayTempRange"`
	NightTempRange string `json:"nightTempRange"`
}
