package task

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store          types.TaskStore
	userStore      types.UserStore
	animalStore    types.AnimalStore
	enclosureStore types.EnclosureStore
}

func NewHandler(store types.TaskStore, userStore types.UserStore, animalStore types.AnimalStore, enclosureStore types.EnclosureStore) *Handler {
	return &Handler{store: store, userStore: userStore, animalStore: animalStore, enclosureStore: enclosureStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleUserCreateTask, h.userStore)).Methods(http.MethodPost)

	// admin routes
	router.HandleFunc("/admin/task", auth.WithJWTAuth(h.handleAdminCreateTask, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleAdminCreateTask(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userId := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userId) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
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

	// check if task exists
	_, err := h.store.GetTaskByNameAndSubjectIdWithUserId(taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId, taskPayload.UserId)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task with name %s and subject ids: %d and %d already exists", taskPayload.TaskName, taskPayload.AnimalId, taskPayload.EnclosureId))
		return
	}

	// if it doesn't exist, create new task
	err = h.store.CreateTask(types.Task{
		TaskName:          taskPayload.TaskName,
		RepeatIntervHours: taskPayload.RepeatIntervHours,
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
	}, taskPayload.AnimalId, taskPayload.EnclosureId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
