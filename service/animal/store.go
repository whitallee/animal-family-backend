package animal

import (
	"database/sql"
	"fmt"

	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateAnimal(animal types.Animal, userID int) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// add animal to animals table
	var addedAnimalId int
	err = tx.QueryRow(`INSERT INTO "animals" ("animalName", "speciesId", "enclosureId", "image", "gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "extraNotes", "isMemorialized", "lastMessage", "memorialPhotos", "memorialDate") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING "animalId"`,
		animal.AnimalName, animal.SpeciesId, animal.EnclosureId, animal.Image, animal.Gender, animal.Dob, animal.PersonalityDesc, animal.DietDesc, animal.RoutineDesc, animal.ExtraNotes, animal.IsMemorialized, animal.LastMessage, animal.MemorialPhotos, animal.MemorialDate).Scan(&addedAnimalId)
	if err != nil {
		return err
	}

	// add user-animal joiner to animalUser table
	_, err = tx.Exec(`INSERT INTO "animalUser" ("animalId", "userId") VALUES ($1, $2)`, addedAnimalId, userID)
	if err != nil {
		return err
	}

	// commit transation
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateAnimal(animal types.Animal) error {
	_, err := s.db.Exec(`UPDATE "animals"
						SET "animalName" = $1, "image" = $2, "extraNotes" = $3, "speciesId" = $4, "enclosureId" = $5,
						"gender" = $6, "dob" = $7, "personalityDesc" = $8, "dietDesc" = $9, "routineDesc" = $10,
						"isMemorialized" = $11, "lastMessage" = $12, "memorialPhotos" = $13, "memorialDate" = $14
						WHERE "animalId" = $15`,
		animal.AnimalName, animal.Image, animal.ExtraNotes, animal.SpeciesId, animal.EnclosureId,
		animal.Gender, animal.Dob, animal.PersonalityDesc, animal.DietDesc, animal.RoutineDesc,
		animal.IsMemorialized, animal.LastMessage, animal.MemorialPhotos, animal.MemorialDate, animal.AnimalId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateAnimalOwner(oldAnimalUser types.AnimalUser, newUserId int) error {
	_, err := s.db.Exec(`UPDATE "animalUser"
						SET "userId" = $1
						WHERE "animalId" = $2 AND "userId" = $3`, newUserId, oldAnimalUser.AnimalId, oldAnimalUser.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAnimals() ([]*types.Animal, error) {
	rows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "isMemorialized", "lastMessage", "memorialPhotos", "memorialDate"
							FROM "animals"`)
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		a, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, a)
	}

	return animals, nil
}

func (s *Store) GetAnimalByNameAndSpeciesWithUserId(animalName string, speciesId int, userID int) (*types.Animal, error) {
	rows, err := s.db.Query(`SELECT a."animalId", a."animalName", a."image", a."extraNotes", a."speciesId", a."enclosureId",
							a."gender", a."dob", a."personalityDesc", a."dietDesc", a."routineDesc", a."isMemorialized", a."lastMessage", a."memorialPhotos", a."memorialDate"
							FROM "animals" a JOIN "animalUser" ON "animalUser"."animalId"=a."animalId"
							WHERE "animalName" = $1 AND "speciesId" = $2 AND "userId" = $3`, animalName, speciesId, userID)
	if err != nil {
		return nil, err
	}

	animal := new(types.Animal)
	for rows.Next() {
		animal, err = utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}
	}

	if animal.AnimalId == 0 {
		return nil, fmt.Errorf("animal not found")
	}

	return animal, nil
}

func (s *Store) GetAnimalUserByIds(animalId int, userID int) (*types.AnimalUser, error) {
	rows, err := s.db.Query(`SELECT * FROM "animalUser" WHERE "animalId" = $1 AND "userId" = $2`, animalId, userID)
	if err != nil {
		return nil, err
	}

	animalUser := new(types.AnimalUser)
	for rows.Next() {
		animalUser, err = utils.ScanRowsIntoAnimalUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if animalUser.AnimalId == 0 && animalUser.UserID == 0 {
		return nil, fmt.Errorf("no ownership found between user and animal")
	}

	return animalUser, nil
}

func (s *Store) GetAnimalUserByAnimalId(animalId int) (*types.AnimalUser, error) {
	rows, err := s.db.Query(`SELECT * FROM "animalUser" WHERE "animalId" = $1`, animalId)
	if err != nil {
		return nil, err
	}

	animalUser := new(types.AnimalUser)
	for rows.Next() {
		animalUser, err = utils.ScanRowsIntoAnimalUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if animalUser.AnimalId == 0 && animalUser.UserID == 0 {
		return nil, fmt.Errorf("no owner of animal with id %d found", animalId)
	}

	return animalUser, nil
}

func (s *Store) GetAnimalById(animalId int) (*types.Animal, error) {
	rows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "isMemorialized", "lastMessage", "memorialPhotos", "memorialDate"
							FROM "animals" WHERE "animalId" = $1`, animalId)
	if err != nil {
		return nil, err
	}

	animal := new(types.Animal)
	for rows.Next() {
		animal, err = utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}
	}

	if animal.AnimalId == 0 {
		return nil, fmt.Errorf("animal not found")
	}

	return animal, nil
}

func (s *Store) GetAnimalsByUserId(userID int) ([]*types.Animal, error) {
	rows, err := s.db.Query(`SELECT a."animalId", a."animalName", a."image", a."extraNotes", a."speciesId", a."enclosureId",
							a."gender", a."dob", a."personalityDesc", a."dietDesc", a."routineDesc", a."isMemorialized", a."lastMessage", a."memorialPhotos", a."memorialDate"
							FROM "animals" a JOIN "animalUser" ON "animalUser"."animalId"=a."animalId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, animal)
	}

	return animals, nil
}

func (s *Store) GetAnimalsByEnclosureId(enclosureId int) ([]*types.Animal, error) {
	rows, err := s.db.Query(`SELECT "animalId", "animalName", "image", "extraNotes", "speciesId", "enclosureId",
							"gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "isMemorialized", "lastMessage", "memorialPhotos", "memorialDate"
							FROM "animals" WHERE "enclosureId" = $1`, enclosureId)
	if err != nil {
		return nil, err
	}

	animals := make([]*types.Animal, 0)
	for rows.Next() {
		animal, err := utils.ScanRowsIntoAnimals(rows)
		if err != nil {
			return nil, err
		}

		animals = append(animals, animal)
	}

	return animals, nil
}

func (s *Store) DeleteAnimalById(animalId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "animalUser" WHERE "animalId" = $1`, animalId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "animals" WHERE "animalId" = $1`, animalId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteAnimalAndTasksById(animalId int) error {
	// get all tasks associated with this animal
	rows, err := s.db.Query(`SELECT "taskId" FROM "taskSubject" WHERE "animalId" = $1`, animalId)
	if err != nil {
		return err
	}

	taskIds := make([]int, 0)
	for rows.Next() {
		var taskId int
		err = rows.Scan(&taskId)
		if err != nil {
			return err
		}
		taskIds = append(taskIds, taskId)
	}
	rows.Close()

	// start deletion transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// delete tasks and their related records
	for _, taskId := range taskIds {
		_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// delete from animalUser and animals
	_, err = tx.Exec(`DELETE FROM "animalUser" WHERE "animalId" = $1`, animalId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM "animals" WHERE "animalId" = $1`, animalId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
