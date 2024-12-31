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
type HabitatStore interface {
	GetHabitats() ([]*Habitat, error)
	CreateHabitat(Habitat) error
}

type Habitat struct {
	HabitatId      int    `json:"habitatId"`
	HabitatName    string `json:"habitatName"`
	HabitatDesc    string `json:"habitatDesc"`
	Image          string `json:"image"`
	Humidity       string `json:"humidity"`
	DayTempRange   string `json:"dayTempRange"`
	NightTempRange string `json:"nightTempRange"`
}

type CreateHabitatPayload struct {
	HabitatName    string `json:"habitatName"`
	HabitatDesc    string `json:"habitatDesc"`
	Image          string `json:"image"`
	Humidity       string `json:"humidity"`
	DayTempRange   string `json:"dayTempRange"`
	NightTempRange string `json:"nightTempRange"`
}

// Enclosure-related Types
type EnclosureStore interface {
	GetEnclosures() ([]*Enclosure, error)
	CreateEnclosure(Enclosure) error
}

type Enclosure struct {
	EnclosureId   int    `json:"enclosureId"`
	EnclosureName string `json:"enclosureName"`
	HabitatId     int    `json:"habitatId"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
}

type CreateEnclosurePayload struct {
	EnclosureName string `json:"enclosureName"`
	HabitatId     int    `json:"habitatId"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
}

// Animal-related Types
type AnimalStore interface {
	GetAnimals() ([]*Animal, error)
	CreateAnimal(Animal) error
	// CreateAnimalWithEnclosure(animal Animal, enclosureId int) error
}

type Animal struct {
	AnimalId    int    `json:"animalId"`
	AnimalName  string `json:"animalName"`
	SpeciesId   int    `json:"speciesId"`
	EnclosureId *int   `json:"enclosureId"`
	Image       string `json:"image"`
	Notes       string `json:"notes"`
}

type CreateAnimalPayload struct {
	AnimalName  string `json:"animalName"`
	SpeciesId   int    `json:"speciesId"`
	EnclosureId *int   `json:"enclosureId"`
	Image       string `json:"image"`
	Notes       string `json:"notes"`
}

// type CreateAnimalWithEnclosurePayload struct {
// 	AnimalName  string `json:"animalName"`
// 	SpeciesId   string `json:"speciesId"`
// 	EnclosureId string `json:"enclosureId"`
// 	Image       string `json:"image"`
// 	Notes       string `json:"notes"`
// }
