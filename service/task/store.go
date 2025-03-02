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

func (s *Store) CreateTask(task types.Task, animalId int, enclosureId int, userId int) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// create task in tasks table
	_, err = tx.Exec("INSERT INTO tasks (taskName, complete, lastCompleted, repeatIntervHours) VALUES (?, ?, ?, ?)",
		task.TaskName, task.Complete, task.LastCompleted, task.RepeatIntervHours)
	if err != nil {
		return err
	}

	// get task id of the newly added task
	var addedTaskId int
	if err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&addedTaskId); err != nil {
		return err
	}

	// add user-task joiner to taskUser table
	_, err = tx.Exec("INSERT INTO taskUser (taskId, userID) VALUES (?,?)", addedTaskId, userId)
	if err != nil {
		return err
	}

	// add subject-task joiner to taskSubject table
	_, err = tx.Exec("INSERT INTO taskSubject (taskId, animalId, enclosureId) VALUES (?,?,?)", addedTaskId, animalId, enclosureId)
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

func (s *Store) UpdateTask(task types.Task) error {
	_, err := s.db.Exec(`UPDATE tasks
						SET taskName = ?, complete = ?, lastCompleted = ?, repeatIntervHours = ?
						WHERE taskId = ?`, task.TaskName, task.Complete, task.LastCompleted, task.RepeatIntervHours, task.TaskId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTaskByNameAndSubjectIdWithUserId(taskName string, animalId int, enclosureId int, userId int) (*types.Task, error) {
	rows, err := s.db.Query(`SELECT t.taskId, t.taskName, t.complete, t.lastCompleted, t.repeatIntervHours
							FROM tasks t JOIN taskUser ON taskUser.taskId=t.taskId
							WHERE taskName = ? AND userId = ?`, taskName, userId) // ADD animalId and enclosureId to the query so it only gets real duplicate tasks
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
	row := s.db.QueryRow("SELECT * FROM taskUser WHERE taskId = ? AND userID = ?", taskId, userID)
	taskUser := new(types.TaskUser)
	err := row.Scan(&taskUser.TaskId, &taskUser.UserID)
	if err != nil {
		return nil, err
	}
	return taskUser, nil
}

func (s *Store) GetTaskById(taskId int) (*types.Task, error) {
	row := s.db.QueryRow("SELECT * FROM tasks WHERE taskId = ?", taskId)
	task := new(types.Task)
	err := row.Scan(&task.TaskId, &task.TaskName, &task.Complete, &task.LastCompleted, &task.RepeatIntervHours)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Store) GetTasksByUserId(userID int) ([]*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE userID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task := new(types.Task)
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Complete, &task.LastCompleted, &task.RepeatIntervHours)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *Store) GetTasksBySubjectIdAndType(subjectId int, taskType string) ([]*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE subjectId = ? AND taskType = ?", subjectId, taskType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		task := new(types.Task)
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Complete, &task.LastCompleted, &task.RepeatIntervHours)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *Store) DeleteTaskById(taskId int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE taskId = ?", taskId)
	return err
}
