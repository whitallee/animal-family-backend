package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
}

type SpeciesStore interface {
	GetSpecies() ([]Species, error)
}

type Species struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	SciName     string `json:"sciName"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Habitat     string `json:"habitat"`
	Diet        string `json:"diet"`
	Sociality   string `json:"sociality"`
	ExtraCare   string `json:"extraCare"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
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
