package task

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

func (s *Store) CheckTaskCompletion() error {
	// check if any tasks should be reset
	_, err := s.db.Exec(`
		UPDATE "tasks" 
		SET "complete" = false 
		WHERE "complete" = true 
		AND "lastCompleted" + ("repeatIntervHours" * interval '1 hour') < NOW()`)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateTask(task types.Task, animalId int, enclosureId int, userId int) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// create task in tasks table
	var addedTaskId int
	err = tx.QueryRow(`INSERT INTO "tasks" ("taskName", "taskDesc", "complete", "lastCompleted", "repeatIntervHours") VALUES ($1, $2, $3, $4, $5) RETURNING "taskId"`,
		task.TaskName, task.TaskDesc, task.Complete, task.LastCompleted, task.RepeatIntervHours).Scan(&addedTaskId)
	if err != nil {
		return err
	}

	// add user-task joiner to taskUser table
	_, err = tx.Exec(`INSERT INTO "taskUser" ("taskId", "userId") VALUES ($1, $2)`, addedTaskId, userId)
	if err != nil {
		return err
	}

	// add subject-task joiner to taskSubject table
	if animalId != 0 && enclosureId == 0 {
		_, err = tx.Exec(`INSERT INTO "taskSubject" ("taskId", "animalId", "enclosureId") VALUES ($1, $2, $3)`, addedTaskId, animalId, nil)
		if err != nil {
			return err
		}
	} else if enclosureId != 0 && animalId == 0 {
		_, err = tx.Exec(`INSERT INTO "taskSubject" ("taskId", "animalId", "enclosureId") VALUES ($1, $2, $3)`, addedTaskId, nil, enclosureId)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid payload, exclusively either animalId or enclosureId must be nonzero")
	}

	// commit transation
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateTask(task types.Task) error {
	_, err := s.db.Exec(`UPDATE "tasks"
						SET "taskName" = $1, "taskDesc" = $2, "complete" = $3, "lastCompleted" = $4, "repeatIntervHours" = $5
						WHERE "taskId" = $6`, task.TaskName, task.TaskDesc, task.Complete, task.LastCompleted, task.RepeatIntervHours, task.TaskId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateTaskOwner(oldTaskUser types.TaskUser, newUserId int) error {
	_, err := s.db.Exec(`UPDATE "taskUser"
						SET "userId" = $1
						WHERE "taskId" = $2 AND "userId" = $3`, newUserId, oldTaskUser.TaskId, oldTaskUser.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateTaskSubject(taskSubject types.TaskSubject) error {
	_, err := s.db.Exec(`UPDATE "taskSubject"
						SET "animalId" = $1, "enclosureId" = $2
						WHERE "taskId" = $3`, taskSubject.AnimalId, taskSubject.EnclosureId, taskSubject.TaskId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTaskByNameAndSubjectIdWithUserId(taskName string, animalId int, enclosureId int, userId int) (*types.Task, error) {
	rows, err := s.db.Query(`SELECT t."taskId", t."taskName", t."taskDesc", t."complete", t."lastCompleted", t."repeatIntervHours"
							FROM "tasks" t JOIN "taskUser" ON "taskUser"."taskId"=t."taskId" JOIN "taskSubject" ON "taskSubject"."taskId"=t."taskId"
							WHERE "taskName" = $1 AND "userId" = $2`, taskName, animalId, enclosureId, userId)
	if err != nil {
		return nil, err
	}

	task := new(types.Task)
	for rows.Next() {
		task, err = utils.ScanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}

	if task.TaskId == 0 {
		return nil, fmt.Errorf("taskId not found")
	}

	return task, nil
}

func (s *Store) GetTaskUserByIds(taskId int, userID int) (*types.TaskUser, error) {
	rows, err := s.db.Query(`SELECT * FROM "taskUser" WHERE "taskId" = $1 AND "userId" = $2`, taskId, userID)
	if err != nil {
		return nil, err
	}

	taskUser := new(types.TaskUser)
	for rows.Next() {
		taskUser, err = utils.ScanRowsIntoTaskUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if taskUser.TaskId == 0 && taskUser.UserID == 0 {
		return nil, fmt.Errorf("no ownership found between user and task")
	}

	return taskUser, nil
}

func (s *Store) GetTaskById(taskId int) (*types.Task, error) {
	rows, err := s.db.Query(`SELECT * FROM "tasks" WHERE "taskId" = $1`, taskId)
	if err != nil {
		return nil, err
	}

	task := new(types.Task)
	for rows.Next() {
		task, err = utils.ScanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}

	if task.TaskId == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

func (s *Store) GetTasksWithSubjectByUserId(userID int) ([]*types.TaskWithSubject, error) {
	rows, err := s.db.Query(`SELECT t."taskId", t."taskName", t."taskDesc", t."complete", t."lastCompleted", t."repeatIntervHours", ts."animalId", ts."enclosureId"
							FROM "tasks" t INNER JOIN "taskUser" tu ON tu."taskId"=t."taskId" INNER JOIN "taskSubject" ts ON ts."taskId"=t."taskId"
							WHERE "userId" = $1`, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.TaskWithSubject, 0)
	for rows.Next() {
		task, err := utils.ScanRowsIntoTaskWithSubject(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Store) GetTasksBySubjectIds(animalId int, enclosureId int) ([]*types.Task, error) {
	rows, err := s.db.Query(`SELECT t."taskId", t."taskName", t."taskDesc", t."complete", t."lastCompleted", t."repeatIntervHours"
							FROM "tasks" t JOIN "taskSubject" ON "taskSubject"."taskId"=t."taskId"
							WHERE "animalId" = $1 AND "enclosureId" = $2`, animalId, enclosureId)
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task := new(types.Task)
		task, err := utils.ScanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Store) DeleteTaskById(taskId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "taskUser" WHERE "taskId" = $1`, taskId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "taskSubject" WHERE "taskId" = $1`, taskId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "tasks" WHERE "taskId" = $1`, taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
