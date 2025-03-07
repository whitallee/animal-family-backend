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

func (s *Store) UpdateTaskSubject(taskSubject types.TaskSubject) error {
	_, err := s.db.Exec(`UPDATE taskSubject
						SET animalId = ?, enclosureId = ?
						WHERE taskId = ?`, taskSubject.AnimalId, taskSubject.EnclosureId, taskSubject.TaskId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTaskByNameAndSubjectIdWithUserId(taskName string, animalId int, enclosureId int, userId int) (*types.Task, error) {
	rows, err := s.db.Query(`SELECT t.taskId, t.taskName, t.complete, t.lastCompleted, t.repeatIntervHours
							FROM tasks t JOIN taskUser ON taskUser.taskId=t.taskId JOIN taskSubject ON taskSubject.taskId=t.taskId
							WHERE taskName = ? AND userId = ?`, taskName, animalId, enclosureId, userId)
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
	rows, err := s.db.Query("SELECT * FROM taskUser WHERE taskId = ? AND userID = ?", taskId, userID)
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
	rows, err := s.db.Query("SELECT * FROM tasks WHERE taskId = ?", taskId)
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

func (s *Store) GetTasksByUserId(userID int) ([]*types.Task, error) {
	rows, err := s.db.Query(`SELECT t.taskId, t.taskName, t.complete, t.lastCompleted, t.repeatIntervHours
							FROM tasks t JOIN taskUser ON taskUser.taskId=t.taskId
							WHERE userId = ?`, userID)
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

func (s *Store) GetTasksBySubjectIds(animalId int, enclosureId int) ([]*types.Task, error) {
	rows, err := s.db.Query(`SELECT t.taskId, t.taskName, t.complete, t.lastCompleted, t.repeatIntervHours
							FROM tasks t JOIN taskSubject ON taskSubject.taskId=t.taskId
							WHERE animalId = ? AND enclosureId = ?`, animalId, enclosureId)
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

	_, err = tx.Exec("DELETE FROM taskUser WHERE taskId = ?", taskId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM taskSubject WHERE taskId = ?", taskId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE taskId = ?", taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
