package task

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// RegisterV2Routes mounts the v2 task routes.
//
// Eight v1 routes collapse into six. GET /task/bysubject becomes a filter on
// the collection, and PUT /task/subject folds into the general update, which
// also subsumes the separate mark-complete and mark-incomplete calls.
//
// Unlike v1, creating or re-pointing a task verifies that the caller owns the
// animal or enclosure it is attached to.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	owned := func(next http.HandlerFunc) http.HandlerFunc {
		return auth.WithJWTAuth(auth.RequireOwnership("id", h.store.UserOwnsTask, next), h.userStore)
	}

	router.HandleFunc("/tasks/check-completion", h.handleCheckTaskCompletionV2).Methods(http.MethodGet)

	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleListTasks, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleCreateTaskV2, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", owned(h.handleGetTask)).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", owned(h.handleUpdateTaskV2)).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", owned(h.handleDeleteTaskV2)).Methods(http.MethodDelete)
}

// handleCheckTaskCompletionV2 godoc
//
//	@Id				checkTaskCompletion
//	@Summary		Reset repeating tasks that are due
//	@Description	Intended for a scheduler. Resets tasks whose repeat interval has elapsed and sends their notifications in the background.
//	@Tags			tasks
//	@Produce		json
//	@Success		200	{object}	types.TaskCompletionResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/tasks/check-completion [get]
func (h *Handler) handleCheckTaskCompletionV2(w http.ResponseWriter, r *http.Request) {
	resetTasks, err := h.store.CheckAndResetTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if len(resetTasks) > 0 && h.notificationSender != nil {
		go h.notificationSender.SendTaskResetNotifications(resetTasks)
	}

	utils.WriteJSON(w, http.StatusOK, types.TaskCompletionResponse{
		TasksReset: len(resetTasks),
	})
}

// handleListTasks godoc
//
//	@Id				listTasks
//	@Summary		List the caller's tasks
//	@Description	Pass animalId or enclosureId to return only the tasks attached to that subject. Supplying both is rejected.
//	@Tags			tasks
//	@Produce		json
//	@Param			animalId	query	int	false	"Only tasks attached to this animal"
//	@Param			enclosureId	query	int	false	"Only tasks attached to this enclosure"
//	@Success		200	{array}		types.TaskWithSubject
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/tasks [get]
func (h *Handler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	animalId, enclosureId, err := parseSubjectFilter(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Always sourced from the caller's own tasks, so the filter narrows a set
	// that is already owner-scoped and cannot expose anyone else's task.
	tasks, err := h.store.GetTasksWithSubjectByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, filterTasksBySubject(tasks, animalId, enclosureId))
}

// parseSubjectFilter reads the optional animalId / enclosureId query
// parameters that replace v1's /task/bysubject route.
func parseSubjectFilter(r *http.Request) (animalId *int, enclosureId *int, err error) {
	animalId, err = optionalPositiveQueryParam(r, "animalId")
	if err != nil {
		return nil, nil, err
	}

	enclosureId, err = optionalPositiveQueryParam(r, "enclosureId")
	if err != nil {
		return nil, nil, err
	}

	if animalId != nil && enclosureId != nil {
		return nil, nil, fmt.Errorf("supply animalId or enclosureId, not both")
	}

	return animalId, enclosureId, nil
}

func optionalPositiveQueryParam(r *http.Request, name string) (*int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return nil, fmt.Errorf("invalid %s: must be a positive integer", name)
	}

	return &value, nil
}

// filterTasksBySubject narrows tasks to one subject. A nil filter returns
// everything.
func filterTasksBySubject(tasks []*types.TaskWithSubject, animalId *int, enclosureId *int) []*types.TaskWithSubject {
	if animalId == nil && enclosureId == nil {
		return tasks
	}

	// Never nil: a nil slice encodes as `null`, but the endpoint is typed as
	// returning an array.
	filtered := make([]*types.TaskWithSubject, 0, len(tasks))
	for _, task := range tasks {
		switch {
		case animalId != nil && task.AnimalId != nil && *task.AnimalId == *animalId:
			filtered = append(filtered, task)
		case enclosureId != nil && task.EnclosureId != nil && *task.EnclosureId == *enclosureId:
			filtered = append(filtered, task)
		}
	}

	return filtered
}

// handleGetTask godoc
//
//	@Id				getTask
//	@Summary		Get one of the caller's tasks
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	types.Task
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/tasks/{id} [get]
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	// Already parsed and ownership-checked by RequireOwnership.
	id := auth.ResourceIDFromContext(r.Context())

	task, err := h.store.GetTaskById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

// handleCreateTaskV2 godoc
//
//	@Id				createTask
//	@Summary		Create a task
//	@Description	A task belongs to exactly one subject, so supply either animalId or enclosureId. The subject must belong to the caller.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body	types.CreateTaskV2Payload	true	"Task to create"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/tasks [post]
func (h *Handler) handleCreateTaskV2(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.CreateTaskV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if err := exactlyOneSubject(payload.AnimalId, payload.EnclosureId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !h.assertOwnsSubject(w, payload.AnimalId, payload.EnclosureId, userID) {
		return
	}

	// The store still takes the zero-as-absent form.
	animalId, enclosureId := zeroIfNil(payload.AnimalId), zeroIfNil(payload.EnclosureId)

	if _, err := h.store.GetTaskByNameAndSubjectIdWithUserId(payload.TaskName, animalId, enclosureId, userID); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task named %s already exists for that subject", payload.TaskName))
		return
	}

	err := h.store.CreateTask(types.Task{
		TaskName:          payload.TaskName,
		TaskDesc:          payload.TaskDesc,
		RepeatIntervHours: payload.RepeatIntervHours,
	}, animalId, enclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusCreated)
}

// handleUpdateTaskV2 godoc
//
//	@Id				updateTask
//	@Summary		Update one of the caller's tasks
//	@Description	A full replace, also used to mark a task complete or incomplete. Supply the subject (animalId or enclosureId) on every update, not only when moving the task.
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Task ID"
//	@Param			task	body	types.UpdateTaskV2Payload	true	"Updated task fields"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/tasks/{id} [put]
func (h *Handler) handleUpdateTaskV2(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.UpdateTaskV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	// The subject is required on every update, not only when it is changing.
	// A task always has exactly one, so treating an omitted subject as "leave
	// it alone" would be the same implicit-preserve behaviour PUT is meant to
	// avoid. Omitting it fails loudly here rather than quietly doing something
	// different from what the request said.
	if err := exactlyOneSubject(payload.AnimalId, payload.EnclosureId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !h.assertOwnsSubject(w, payload.AnimalId, payload.EnclosureId, userID) {
		return
	}

	err := h.store.UpdateTask(types.Task{
		TaskId:            id,
		TaskName:          payload.TaskName,
		TaskDesc:          payload.TaskDesc,
		Complete:          payload.Complete,
		LastCompleted:     payload.LastCompleted,
		RepeatIntervHours: payload.RepeatIntervHours,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.SetTaskSubject(id, payload.AnimalId, payload.EnclosureId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// handleDeleteTaskV2 godoc
//
//	@Id				deleteTask
//	@Summary		Delete one of the caller's tasks
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/tasks/{id} [delete]
func (h *Handler) handleDeleteTaskV2(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	if err := h.store.DeleteTaskById(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// exactlyOneSubject enforces the rule the schema implies: "taskSubject" holds
// two nullable foreign keys and a task belongs to one subject.
func exactlyOneSubject(animalId *int, enclosureId *int) error {
	switch {
	case animalId == nil && enclosureId == nil:
		return fmt.Errorf("a task needs a subject: supply animalId or enclosureId")
	case animalId != nil && enclosureId != nil:
		return fmt.Errorf("a task has one subject: supply animalId or enclosureId, not both")
	default:
		return nil
	}
}

// assertOwnsSubject verifies the caller owns whichever subject was supplied,
// writing the response and returning false if not.
//
// v1 omitted this check entirely, so a task could be attached to another user's
// animal or enclosure.
func (h *Handler) assertOwnsSubject(w http.ResponseWriter, animalId *int, enclosureId *int, userID int) bool {
	var (
		owned bool
		err   error
	)

	switch {
	case animalId != nil:
		owned, err = h.animalStore.UserOwnsAnimal(*animalId, userID)
	case enclosureId != nil:
		owned, err = h.enclosureStore.UserOwnsEnclosure(*enclosureId, userID)
	default:
		return true
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return false
	}

	if !owned {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you do not have access to that subject"))
		return false
	}

	return true
}

// zeroIfNil converts to the zero-as-absent form the existing store expects.
func zeroIfNil(value *int) int {
	if value == nil {
		return 0
	}

	return *value
}
