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
	UpdateSpecies(Species) error // TODO
	GetSpecies() ([]*Species, error)
	GetSpeciesByComName(string) (*Species, error)
	GetSpeciesBySciName(string) (*Species, error)
	GetSpeciesById(int) (*Species, error) // not used in any handler yet
	DeleteSpeciesById(int) error
}

type Species struct {
	SpeciesID          int    `json:"speciesId"`
	ComName            string `json:"comName"`
	SciName            string `json:"sciName"`
	SpeciesDesc        string `json:"speciesDesc"`
	Image              string `json:"image"`
	HabitatId          int    `json:"habitatId"`
	BaskTemp           string `json:"baskTemp"`
	Diet               string `json:"diet"`
	Sociality          string `json:"sociality"`
	Lifespan           string `json:"lifespan"`
	Size               string `json:"size"`
	Weight             string `json:"weight"`
	ConservationStatus string `json:"conservationStatus"`
	ExtraCare          string `json:"extraCare"`
}

type CreateSpeciesPayload struct {
	ComName            string `json:"comName" validate:"required"`
	SciName            string `json:"sciName" validate:"required"`
	SpeciesDesc        string `json:"speciesDesc" validate:"required"`
	Image              string `json:"image" validate:"required"`
	HabitatId          int    `json:"habitatId" validate:"required,min=0"`
	BaskTemp           string `json:"baskTemp" validate:"required"`
	Diet               string `json:"diet" validate:"required"`
	Sociality          string `json:"sociality" validate:"required"`
	Lifespan           string `json:"lifespan" validate:"required"`
	Size               string `json:"size" validate:"required"`
	Weight             string `json:"weight" validate:"required"`
	ConservationStatus string `json:"conservationStatus" validate:"required"`
	ExtraCare          string `json:"extraCare" validate:"required"`
}

type UpdateSpeciesPayload struct {
	SpeciesID          int    `json:"speciesId" validate:"required"`
	ComName            string `json:"comName" validate:"required"`
	SciName            string `json:"sciName" validate:"required"`
	SpeciesDesc        string `json:"speciesDesc" validate:"required"`
	Image              string `json:"image" validate:"required"`
	HabitatId          int    `json:"habitatId" validate:"required,min=0"`
	BaskTemp           string `json:"baskTemp" validate:"required"`
	Diet               string `json:"diet" validate:"required"`
	Sociality          string `json:"sociality" validate:"required"`
	Lifespan           string `json:"lifespan" validate:"required"`
	Size               string `json:"size" validate:"required"`
	Weight             string `json:"weight" validate:"required"`
	ConservationStatus string `json:"conservationStatus" validate:"required"`
	ExtraCare          string `json:"extraCare" validate:"required"`
}

type SpeciesIdPayload struct {
	SpeciesId int `json:"speciesId" validate:"required,min=0"`
}

// Habitat-related types
type HabitatStore interface {
	CreateHabitat(Habitat) error
	UpdateHabitat(Habitat) error
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

type UpdateHabitatPayload struct {
	HabitatId      int    `json:"habitatId" validate:"required"`
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
	CreateEnclosure(Enclosure, int) error
	CreateEnclosureWithAnimals(enclosure Enclosure, animalIds []int, userID int) error
	UpdateEnclosure(Enclosure) error
	UpdateEnclosureOwnerWithAnimals(oldEnclosureUser EnclosureUser, newUserId int) error
	GetEnclosures() ([]*Enclosure, error)
	GetEnclosureByNameAndHabitatWithUserId(enclosureName string, habitatId int, userID int) (*Enclosure, error)
	GetEnclosureUserByIds(enclosureId int, userID int) (*EnclosureUser, error)
	GetEnclosureUserByEnclosureId(enclosureId int) (*EnclosureUser, error)
	GetEnclosuresByUserId(int) ([]*Enclosure, error)
	GetEnclosureById(int) (*Enclosure, error)
	DeleteEnclosureById(enclosureId int) error
	DeleteEnclosureAndAnimalsById(enclosureId int) error
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

type CreateEnclosureWithOwnerPayload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	UserID        int    `json:"userId" validate:"required,min=0"`
}

type CreateEnclosureWithAnimalsPayload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	AnimalIds     []int  `json:"animalIds" validate:"required"`
}

type CreateEnclosureWithOwnerWithAnimalsPayload struct {
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	AnimalIds     []int  `json:"animalIds" validate:"required"`
	UserID        int    `json:"userId" validate:"required,min=0"`
}

type UpdateEnclosurePayload struct {
	EnclosureId   int    `json:"enclosureId" validate:"required"`
	EnclosureName string `json:"enclosureName" validate:"required"`
	HabitatId     int    `json:"habitatId" validate:"required,min=0"`
	Image         string `json:"image" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
}

type UpdateEnclosureOwnerPayload struct {
	EnclosureId int `json:"enclosureId" validate:"required,min=0"`
	OldUserId   int `json:"oldUserId" validate:"required,min=0"`
	NewUserId   int `json:"newUserId" validate:"required,min=0"`
}

type EnclosureIdPayload struct {
	EnclosureId int `json:"enclosureId" validate:"required,min=0"`
}

// Animal-related Types
type AnimalStore interface {
	CreateAnimal(Animal, int) error
	UpdateAnimal(Animal) error
	UpdateAnimalOwner(oldAnimalUser AnimalUser, newUserId int) error
	GetAnimals() ([]*Animal, error)
	GetAnimalByNameAndSpeciesWithUserId(animalName string, speciesId int, userID int) (*Animal, error)
	GetAnimalUserByIds(animalId int, userID int) (*AnimalUser, error)
	GetAnimalUserByAnimalId(animalId int) (*AnimalUser, error)
	GetAnimalById(int) (*Animal, error)
	GetAnimalsByUserId(int) ([]*Animal, error)
	GetAnimalsByEnclosureId(int) ([]*Animal, error)
	DeleteAnimalById(int) error
}

type Animal struct {
	AnimalId        int       `json:"animalId"`
	AnimalName      string    `json:"animalName"`
	SpeciesId       int       `json:"speciesId"`
	EnclosureId     *int      `json:"enclosureId"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender"`
	Dob             time.Time `json:"dob"`
	PersonalityDesc string    `json:"personalityDesc"`
	DietDesc        string    `json:"dietDesc"`
	RoutineDesc     string    `json:"routineDesc"`
	ExtraNotes      string    `json:"extraNotes"`
}

type AnimalUser struct {
	AnimalId int `json:"animalId"`
	UserID   int `json:"userID"`
}

type CreateAnimalWithOwnerPayload struct {
	AnimalName      string    `json:"animalName" validate:"required"`
	SpeciesId       int       `json:"speciesId" validate:"required,min=0"`
	EnclosureId     *int      `json:"enclosureId" validate:"required,min=0"`
	Image           string    `json:"image" validate:"required"`
	Gender          string    `json:"gender" validate:"required"`
	Dob             time.Time `json:"dob" validate:"required"`
	PersonalityDesc string    `json:"personalityDesc" validate:"required"`
	DietDesc        string    `json:"dietDesc" validate:"required"`
	RoutineDesc     string    `json:"routineDesc" validate:"required"`
	ExtraNotes      string    `json:"extraNotes" validate:"required"`
	UserID          int       `json:"userId" validate:"required,min=0"`
}

type CreateAnimalPayload struct {
	AnimalName      string    `json:"animalName" validate:"required"`
	SpeciesId       int       `json:"speciesId" validate:"required,min=0"`
	EnclosureId     *int      `json:"enclosureId" validate:"omitempty,min=0"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender"`
	Dob             time.Time `json:"dob"`
	PersonalityDesc string    `json:"personalityDesc"`
	DietDesc        string    `json:"dietDesc"`
	RoutineDesc     string    `json:"routineDesc"`
	ExtraNotes      string    `json:"extraNotes"`
}

type UpdateAnimalPayload struct {
	AnimalId        int       `json:"animalId" validate:"required"`
	AnimalName      string    `json:"animalName" validate:"required"`
	SpeciesId       int       `json:"speciesId" validate:"required,min=0"`
	EnclosureId     *int      `json:"enclosureId" validate:"required,min=0"`
	Image           string    `json:"image" validate:"required"`
	Gender          string    `json:"gender" validate:"required"`
	Dob             time.Time `json:"dob" validate:"required"`
	PersonalityDesc string    `json:"personalityDesc" validate:"required"`
	DietDesc        string    `json:"dietDesc" validate:"required"`
	RoutineDesc     string    `json:"routineDesc" validate:"required"`
	ExtraNotes      string    `json:"extraNotes" validate:"required"`
}

type UpdateAnimalOwnerPayload struct {
	AnimalId  int `json:"animalId" validate:"required,min=0"`
	OldUserId int `json:"oldUserId" validate:"required,min=0"`
	NewUserId int `json:"newUserId" validate:"required,min=0"`
}

type AnimalIdPayload struct {
	AnimalId int `json:"animalId" validate:"required,min=0"`
}

// Task-related Types
type TaskStore interface {
	CheckTaskCompletion() error
	CreateTask(task Task, animalId int, enclosureId int, userId int) error
	UpdateTask(Task) error
	UpdateTaskOwner(oldTaskUser TaskUser, newUserId int) error
	UpdateTaskSubject(TaskSubject) error
	GetTaskByNameAndSubjectIdWithUserId(taskName string, animalId int, enclosureId int, userId int) (*Task, error)
	GetTaskUserByIds(taskId int, userID int) (*TaskUser, error)
	GetTaskById(int) (*Task, error)
	GetTasksWithSubjectByUserId(int) ([]*TaskWithSubject, error)
	GetTasksBySubjectIds(animalId int, enclosureId int) ([]*Task, error)
	DeleteTaskById(int) error
}

type Task struct {
	TaskId            int       `json:"taskId"`
	TaskName          string    `json:"taskName"`
	TaskDesc          string    `json:"taskDesc"`
	Complete          bool      `json:"complete"`
	LastCompleted     time.Time `json:"lastCompleted"`
	RepeatIntervHours int       `json:"repeatIntervHours"`
}

type TaskWithSubject struct {
	TaskId            int       `json:"taskId"`
	TaskName          string    `json:"taskName"`
	TaskDesc          string    `json:"taskDesc"`
	Complete          bool      `json:"complete"`
	LastCompleted     time.Time `json:"lastCompleted"`
	RepeatIntervHours int       `json:"repeatIntervHours"`
	AnimalId          *int      `json:"animalId"`
	EnclosureId       *int      `json:"enclosureId"`
}

type TaskUser struct {
	TaskId int `json:"taskId"`
	UserID int `json:"userID"`
}

type TaskSubject struct {
	TaskId      int `json:"taskId"`
	AnimalId    int `json:"animalId"`
	EnclosureId int `json:"enclosureId"`
}

type CreateTaskPayload struct {
	TaskName          string `json:"taskName" validate:"required"`
	TaskDesc          string `json:"taskDesc" validate:"required"`
	RepeatIntervHours int    `json:"repeatIntervHours" validate:"required,min=0"`
	AnimalId          int    `json:"animalId" validate:"min=0"`
	EnclosureId       int    `json:"enclosureId" validate:"min=0"`
}

type CreateTaskWithOwnerPayload struct {
	TaskName          string `json:"taskName" validate:"required"`
	TaskDesc          string `json:"taskDesc" validate:"required"`
	RepeatIntervHours int    `json:"repeatIntervHours" validate:"required,min=1"`
	AnimalId          int    `json:"animalId" validate:"min=0"`
	EnclosureId       int    `json:"enclosureId" validate:"min=0"`
	UserId            int    `json:"userId" validate:"required,min=0"`
}

type UpdateTaskPayload struct {
	TaskId            int       `json:"taskId" validate:"required,min=0"`
	TaskName          string    `json:"taskName" validate:"required"`
	TaskDesc          string    `json:"taskDesc" validate:"required"`
	Complete          bool      `json:"complete" validate:""`
	LastCompleted     time.Time `json:"lastCompleted" validate:"required"`
	RepeatIntervHours int       `json:"repeatIntervHours" validate:"required,min=0"`
}

type UpdateTaskOwnerPayload struct {
	TaskId    int `json:"taskId" validate:"required,min=0"`
	OldUserId int `json:"oldUserId" validate:"required,min=0"`
	NewUserId int `json:"newUserId" validate:"required,min=0"`
}

type UpdateTaskSubjectPayload struct {
	TaskId      int `json:"taskId" validate:"required,min=0"`
	AnimalId    int `json:"animalId" validate:"min=0"`
	EnclosureId int `json:"enclosureId" validate:"min=0"`
}

type TaskIdPayload struct {
	TaskId int `json:"taskId" validate:"required,min=0"`
}

type SubjectIdsPayload struct {
	AnimalId    int `json:"animalId" validate:"min=0"`
	EnclosureId int `json:"enclosureId" validate:"min=0"`
}

type LoopMessageStore interface {
	ReceiveLoopMessage(InboundLoopMessagePayload) error
	SendLoopMessage(SentLoopMessagePayload) error
}

type InboundLoopMessagePayload struct {
	AlertType   string `json:"alertType"`
	Recipient   string `json:"recipient"`
	Text        string `json:"text"`
	MessageType string `json:"messageType"`
	MessageId   string `json:"messageId"`
	WebhookId   string `json:"webhookId"`
	ApiVersion  string `json:"apiVersion"`
}

type SentLoopMessagePayload struct {
	AlertType  string `json:"alertType"`
	Success    bool   `json:"success"`
	Recipient  string `json:"recipient"`
	Text       string `json:"text"`
	MessageId  string `json:"messageId"`
	WebhookId  string `json:"webhookId"`
	ApiVersion string `json:"apiVersion"`
}
