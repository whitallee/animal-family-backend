package task

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/service/notification"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store              types.TaskStore
	userStore          types.UserStore
	animalStore        types.AnimalStore
	enclosureStore     types.EnclosureStore
	notificationSender *notification.NotificationSender
}

func NewHandler(store types.TaskStore, userStore types.UserStore, animalStore types.AnimalStore, enclosureStore types.EnclosureStore, notificationSender *notification.NotificationSender) *Handler {
	return &Handler{store: store, userStore: userStore, animalStore: animalStore, enclosureStore: enclosureStore, notificationSender: notificationSender}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// public routes
	router.HandleFunc("/task/check-completion", h.handleCheckTaskCompletion).Methods(http.MethodGet)

	// user routes
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleUserGetTasks, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/task/byid", auth.WithJWTAuth(h.handleUserGetTaskById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/task/bysubject", auth.WithJWTAuth(h.handleUserGetTasksBySubject, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleUserCreateTask, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleUserUpdateTask, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/task/subject", auth.WithJWTAuth(h.handleUserUpdateTaskSubject, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleUserDeleteTask, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/admin/task/byid", auth.WithJWTAuth(h.handleAdminGetTaskById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/task/byuser", auth.WithJWTAuth(h.handleAdminGetTasksByUser, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/task/bysubject", auth.WithJWTAuth(h.handleAdminGetTasksBySubject, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/task", auth.WithJWTAuth(h.handleAdminCreateTask, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/task", auth.WithJWTAuth(h.handleAdminUpdateTask, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/task/owner", auth.WithJWTAuth(h.handleAdminUpdateTaskOwner, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/task/subject", auth.WithJWTAuth(h.handleAdminUpdateTaskSubject, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/task", auth.WithJWTAuth(h.handleAdminDeleteTask, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleCheckTaskCompletion(w http.ResponseWriter, r *http.Request) {
	// Get tasks that need resetting
	resetTasks, err := h.store.CheckAndResetTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Send notifications asynchronously (non-blocking)
	if len(resetTasks) > 0 && h.notificationSender != nil {
		go h.notificationSender.SendTaskResetNotifications(resetTasks)
	}

	// Return immediately
	utils.WriteJSON(w, http.StatusOK, map[string]int{
		"tasksReset": len(resetTasks),
	})
}

func (h *Handler) handleAdminCreateTask(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var taskPayload types.CreateTaskWithOwnerPayload
	if err := utils.ParseJSON(r, &taskPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check that only 1 subject is non-0
	if (taskPayload.AnimalId == 0 && taskPayload.EnclosureId == 0) || (taskPayload.AnimalId != 0 && taskPayload.EnclosureId != 0) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload, exclusively either animalId or enclosureId must be nonzero"))
		return
	}

	// check if task exists DOESNT PREVENT DUPES!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	_, err := h.store.GetTaskByNameAndSubjectIdWithUserId(taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId, taskPayload.UserId)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task with name %s and subject ids: %d and %d already exists", taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId))
		return
	}

	// if it doesn't exist, create new task
	err = h.store.CreateTask(types.Task{
		TaskName:          taskPayload.TaskName,
		RepeatIntervHours: taskPayload.RepeatIntervHours,
		LastCompleted:     time.Now(),
	}, taskPayload.AnimalId, taskPayload.EnclosureId, taskPayload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUserCreateTask(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var taskPayload types.CreateTaskPayload
	if err := utils.ParseJSON(r, &taskPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check that only 1 subject is non-0
	if (taskPayload.AnimalId == 0 && taskPayload.EnclosureId == 0) || (taskPayload.AnimalId != 0 && taskPayload.EnclosureId != 0) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload, exclusively either animalId or enclosureId must be nonzero"))
		return
	}

	// check if task exists
	_, err := h.store.GetTaskByNameAndSubjectIdWithUserId(taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId, userId)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task with name %s and subject ids: %d and %d already exists", taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId))
		return
	}

	// if it doesn't exist, create new task
	err = h.store.CreateTask(types.Task{
		TaskName:          taskPayload.TaskName,
		RepeatIntervHours: taskPayload.RepeatIntervHours,
		LastCompleted:     time.Now(),
	}, taskPayload.AnimalId, taskPayload.EnclosureId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleAdminUpdateTask(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var taskPayload types.UpdateTaskPayload
	if err := utils.ParseJSON(r, &taskPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update task
	err := h.store.UpdateTask(types.Task(taskPayload))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminUpdateTaskOwner(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var updateTaskOwnerPayload types.UpdateTaskOwnerPayload
	if err := utils.ParseJSON(r, &updateTaskOwnerPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(updateTaskOwnerPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update taskOwner
	err := h.store.UpdateTaskOwner(types.TaskUser{
		TaskId: updateTaskOwnerPayload.TaskId,
		UserID: updateTaskOwnerPayload.OldUserId,
	}, updateTaskOwnerPayload.NewUserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserUpdateTask(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var taskPayload types.UpdateTaskPayload
	if err := utils.ParseJSON(r, &taskPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if ownership exists
	_, err := h.store.GetTaskUserByIds(taskPayload.TaskId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, update task
	err = h.store.UpdateTask(types.Task(taskPayload))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminUpdateTaskSubject(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var taskSubjectPayload types.UpdateTaskSubjectPayload
	if err := utils.ParseJSON(r, &taskSubjectPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskSubjectPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check that only 1 subject is non-0
	if (taskSubjectPayload.AnimalId == 0 && taskSubjectPayload.EnclosureId == 0) || (taskSubjectPayload.AnimalId != 0 && taskSubjectPayload.EnclosureId != 0) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload, exclusively either animalId or enclosureId must be nonzero"))
		return
	}

	// update task
	err := h.store.UpdateTaskSubject(types.TaskSubject(taskSubjectPayload))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserUpdateTaskSubject(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var taskSubjectPayload types.UpdateTaskSubjectPayload
	if err := utils.ParseJSON(r, &taskSubjectPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskSubjectPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check that only 1 subject is non-0
	if (taskSubjectPayload.AnimalId == 0 && taskSubjectPayload.EnclosureId == 0) || (taskSubjectPayload.AnimalId != 0 && taskSubjectPayload.EnclosureId != 0) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload, exclusively either animalId or enclosureId must be nonzero"))
		return
	}

	// check if ownership exists
	_, err := h.store.GetTaskUserByIds(taskSubjectPayload.TaskId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, update task
	err = h.store.UpdateTaskSubject(types.TaskSubject(taskSubjectPayload))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminGetTasksByUser(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var userIdPayload types.UserIDPayload
	if err := utils.ParseJSON(r, &userIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(userIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get tasks
	taskList, err := h.store.GetTasksWithSubjectByUserId(userIdPayload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, taskList)

}

func (h *Handler) handleUserGetTasks(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get tasks
	taskList, err := h.store.GetTasksWithSubjectByUserId(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, taskList)

}

func (h *Handler) handleAdminGetTaskById(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var taskIdPayload types.TaskIdPayload
	if err := utils.ParseJSON(r, &taskIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get task
	task, err := h.store.GetTaskById(taskIdPayload.TaskId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)

}

func (h *Handler) handleUserGetTaskById(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var taskIdPayload types.TaskIdPayload
	if err := utils.ParseJSON(r, &taskIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if ownership exists
	_, err := h.store.GetTaskUserByIds(taskIdPayload.TaskId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// get task
	task, err := h.store.GetTaskById(taskIdPayload.TaskId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)

}

func (h *Handler) handleAdminGetTasksBySubject(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var subjectIdsPayload types.SubjectIdsPayload
	if err := utils.ParseJSON(r, &subjectIdsPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(subjectIdsPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get tasks
	taskList, err := h.store.GetTasksBySubjectIds(subjectIdsPayload.AnimalId, subjectIdsPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, taskList)

}

func (h *Handler) handleUserGetTasksBySubject(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var subjectIdsPayload types.SubjectIdsPayload
	if err := utils.ParseJSON(r, &subjectIdsPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(subjectIdsPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if ownership exists
	if subjectIdsPayload.AnimalId != 0 {
		_, err := h.animalStore.GetAnimalUserByIds(subjectIdsPayload.AnimalId, userId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
			return
		}
	} else if subjectIdsPayload.EnclosureId != 0 {
		_, err := h.enclosureStore.GetEnclosureUserByIds(subjectIdsPayload.EnclosureId, userId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
			return
		}
	}

	// if ownership exists, get tasks
	taskList, err := h.store.GetTasksBySubjectIds(subjectIdsPayload.AnimalId, subjectIdsPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, taskList)

}

func (h *Handler) handleAdminDeleteTask(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
		return
	}

	// get JSON payload
	var taskIdPayload types.TaskIdPayload
	if err := utils.ParseJSON(r, &taskIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete task
	err := h.store.DeleteTaskById(taskIdPayload.TaskId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserDeleteTask(w http.ResponseWriter, r *http.Request) {
	// get userId
	userId := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var taskIdPayload types.TaskIdPayload
	if err := utils.ParseJSON(r, &taskIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(taskIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check for ownership
	_, err := h.store.GetTaskUserByIds(taskIdPayload.TaskId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists delete task
	err = h.store.DeleteTaskById(taskIdPayload.TaskId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
