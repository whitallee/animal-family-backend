package types

import (
	"database/sql"
	"time"
)

// User-related types
type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
}

type User struct {
	ID        int            `json:"id"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Email     string         `json:"email"`
	Phone     sql.NullString `json:"phone"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
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
	CreateSpecies(Species) error
	GetSpecies() ([]*Species, error)
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
	CreateHabitat(Habitat) error
	GetHabitats() ([]*Habitat, error)
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
	CreateEnclosure(Enclosure) error
	CreateEnclosureByUserId(Enclosure, int) error
	CreateEnclosureWithAnimalsByUserId(Enclosure, []int, int) error
	GetEnclosures() ([]*Enclosure, error)
	GetEnclosuresByUserId(int) ([]*Enclosure, error)
	GetEnclosureByIdWithUserId(enclosureId int, userID int) (*Enclosure, error)
	DeleteEnclosureAndAnimalsByIdWithUserId(enclosureId int, userID int) error
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

type CreateEnclosureWithAnimalsPayload struct {
	EnclosureName string `json:"enclosureName"`
	HabitatId     int    `json:"habitatId"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
	AnimalIds     []int  `json:"animalIds"`
}

type GetEnclosureByIdWithUserIdPayload struct {
	EnclosureId int `json:"enclosureId"`
}

type DeleteEnclosureAndAnimalsByUserIdPayload struct {
	EnclosureId int `json:"enclosureId"`
}

// Animal-related Types
type AnimalStore interface {
	CreateAnimal(Animal) error
	CreateAnimalByUserId(Animal, int) error
	GetAnimals() ([]*Animal, error)
	GetAnimalsByUserId(int) ([]*Animal, error)
	GetAnimalsByEnclosureIdWithUserId(enclosureId int, userID int) ([]*Animal, error)
	DeleteAnimalByIdWithUserId(animalId int, userID int) error
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

type GetAnimalsByEnclosureIdWithUserIdPayload struct {
	EnclosureId int `json:"enclosureId"`
}

type DeleteAnimalByIdWithUserIdPayload struct {
	AnimalId int `json:"animalId"`
}
