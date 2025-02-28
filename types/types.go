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
	DeleteUserById(id int) error
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

type UserIDPayload struct {
	UserID int `json:"userId" validate:"required,min=0"`
}

type UserEmailPayload struct {
	UserEmail string `json:"email" validate:"required,email"`
}

// Species-related types
type SpeciesStore interface {
	CreateSpecies(Species) error
	UpdateSpeciesById(int) error // TODO
	GetSpecies() ([]*Species, error)
	GetSpeciesByComName(string) (*Species, error)
	GetSpeciesBySciName(string) (*Species, error)
	GetSpeciesById(int) (*Species, error) // not used in any handler yet
	DeleteSpeciesById(int) error
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
	ComName     string `json:"comName" validate:"required"`
	SciName     string `json:"sciName" validate:"required"`
	SpeciesDesc string `json:"speciesDesc" validate:"required"`
	Image       string `json:"image" validate:"required"`
	HabitatId   int    `json:"habitatId" validate:"required,min=0"`
	BaskTemp    string `json:"baskTemp" validate:"required"`
	Diet        string `json:"diet" validate:"required"`
	Sociality   string `json:"sociality" validate:"required"`
	ExtraCare   string `json:"extraCare" validate:"required"`
}

type SpeciesIdPayload struct {
	SpeciesId int `json:"speciesId" validate:"required,min=0"`
}

// Habitat-related types
type HabitatStore interface {
	CreateHabitat(Habitat) error
	UpdateHabitatById(int) error // TODO
	GetHabitats() ([]*Habitat, error)
	GetHabitatByName(string) (*Habitat, error)
	GetHabitatById(int) (*Habitat, error) // not used in any handler yet
	DeleteHabitatById(int) error
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
	HabitatName    string `json:"habitatName" validate:"required"`
	HabitatDesc    string `json:"habitatDesc" validate:"required"`
	Image          string `json:"image" validate:"required"`
	Humidity       string `json:"humidity" validate:"required"`
	DayTempRange   string `json:"dayTempRange" validate:"required"`
	NightTempRange string `json:"nightTempRange" validate:"required"`
}

type HabitatIdPayload struct {
	HabitatId int `json:"habitatId" validate:"required,min=0"`
}

// Enclosure-related Types
type EnclosureStore interface {
	CreateEnclosure(Enclosure) error
	CreateEnclosureByUserId(Enclosure, int) error
	CreateEnclosureWithAnimalsByUserId(enclosure Enclosure, animalIds []int, userID int) error
	UpdateEnclosure(Enclosure) error
	GetEnclosures() ([]*Enclosure, error)
	GetEnclosureByNameAndHabitatWithUserId(enclosureName string, habitatId int, userID int) (*Enclosure, error)
	GetEnclosureUserByIds(enclosureId int, userID int) (*EnclosureUser, error)
	GetEnclosuresByUserId(int) ([]*Enclosure, error)
	GetEnclosureByIdWithUserId(enclosureId int, userID int) (*Enclosure, error)
	DeleteEnclosureByIdWithUserId(enclosureId int, userID int) error
	DeleteEnclosureAndAnimalsByIdWithUserId(enclosureId int, userID int) error
}

type Enclosure struct {
	EnclosureId   int    `json:"enclosureId"`
	EnclosureName string `json:"enclosureName"`
	HabitatId     int    `json:"habitatId"`
	Image         string `json:"image"`
	Notes         string `json:"notes"`
}

type EnclosureUser struct {
	EnclosureId int `json:"enclosureId"`
	UserID      int `json:"userID"`
}

type CreateEnclosurePayload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
}

type CreateEnclosureWithAnimalsPayload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	AnimalIds     []int  `json:"animalIds" validate:"required"`
}

type UpdateEnclosurePayload struct {
	EnclosureId   int    `json:"enclosureId" validate:"required"`
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
}

type EnclosureIdPayload struct {
	EnclosureId int `json:"enclosureId" validate:"required,min=0"`
}

// Animal-related Types
type AnimalStore interface {
	CreateAnimal(Animal) error
	CreateAnimalByUserId(Animal, int) error
	UpdateAnimal(Animal) error
	GetAnimals() ([]*Animal, error)
	GetAnimalByNameAndSpeciesWithUserId(animalName string, speciesId int, userID int) (*Animal, error)
	GetAnimalUserByIds(animalId int, userID int) (*AnimalUser, error)
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

type AnimalUser struct {
	AnimalId int `json:"animalId"`
	UserID   int `json:"userID"`
}

type CreateAnimalPayload struct {
	AnimalName  string `json:"animalName" validate:"required"`
	SpeciesId   int    `json:"speciesId" validate:"required,min=0"`
	EnclosureId *int   `json:"enclosureId" validate:"required,min=0"`
	Image       string `json:"image" validate:"required"`
	Notes       string `json:"notes" validate:"required"`
}

type UpdateAnimalPayload struct {
	AnimalId    int    `json:"animalId" validate:"required"`
	AnimalName  string `json:"animalName" validate:"required"`
	SpeciesId   int    `json:"speciesId" validate:"required,min=0"`
	EnclosureId *int   `json:"enclosureId" validate:"required,min=0"`
	Image       string `json:"image" validate:"required"`
	Notes       string `json:"notes" validate:"required"`
}

type AnimalIdPayload struct {
	AnimalId int `json:"animalId" validate:"required,min=0"`
}
