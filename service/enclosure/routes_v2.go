package enclosure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// RegisterV2Routes mounts the v2 enclosure routes.
//
// Differences from v1: plural collection path, IDs as path parameters, the
// separate /withanimals create route folded into POST /enclosures via an
// optional animalIds field, and the two cascading delete routes replaced by a
// ?cascade= query parameter. Ownership is enforced by middleware rather than
// repeated in each handler.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	owned := func(next http.HandlerFunc) http.HandlerFunc {
		return auth.WithJWTAuth(auth.RequireOwnership("id", h.store.UserOwnsEnclosure, next), h.userStore)
	}

	router.HandleFunc("/enclosures", auth.WithJWTAuth(h.handleListEnclosures, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/enclosures", auth.WithJWTAuth(h.handleCreateEnclosure, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/enclosures/{id}", owned(h.handleGetEnclosure)).Methods(http.MethodGet)
	router.HandleFunc("/enclosures/{id}", owned(h.handleUpdateEnclosure)).Methods(http.MethodPut)
	router.HandleFunc("/enclosures/{id}", owned(h.handleDeleteEnclosure)).Methods(http.MethodDelete)
}

// handleListEnclosures godoc
//
//	@Id				listEnclosures
//	@Summary		List the caller's enclosures
//	@Tags			enclosures
//	@Produce		json
//	@Success		200	{array}		types.Enclosure
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/enclosures [get]
func (h *Handler) handleListEnclosures(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	enclosures, err := h.store.GetEnclosuresByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosures)
}

// handleGetEnclosure godoc
//
//	@Id				getEnclosure
//	@Summary		Get one of the caller's enclosures
//	@Tags			enclosures
//	@Produce		json
//	@Param			id	path		int	true	"Enclosure ID"
//	@Success		200	{object}	types.Enclosure
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/enclosures/{id} [get]
func (h *Handler) handleGetEnclosure(w http.ResponseWriter, r *http.Request) {
	// Already parsed and ownership-checked by RequireOwnership.
	id := auth.ResourceIDFromContext(r.Context())

	enclosure, err := h.store.GetEnclosureById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosure)
}

// handleCreateEnclosure godoc
//
//	@Id				createEnclosure
//	@Summary		Create an enclosure
//	@Description	Supply animalIds to move existing animals into the new enclosure as it is created. Every animal listed must belong to the caller.
//	@Tags			enclosures
//	@Accept			json
//	@Produce		json
//	@Param			enclosure	body	types.CreateEnclosureV2Payload	true	"Enclosure to create"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/enclosures [post]
func (h *Handler) handleCreateEnclosure(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.CreateEnclosureV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if _, err := h.store.GetEnclosureByNameAndHabitatWithUserId(payload.EnclosureName, payload.HabitatId, userID); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s already exists in that habitat", payload.EnclosureName))
		return
	}

	enclosure := types.Enclosure{
		EnclosureName: payload.EnclosureName,
		HabitatId:     payload.HabitatId,
		Image:         payload.Image,
		Notes:         payload.Notes,
	}

	if len(payload.AnimalIds) == 0 {
		if err := h.store.CreateEnclosure(enclosure, userID); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, nil)
		return
	}

	if err := h.store.CreateEnclosureWithAnimals(enclosure, payload.AnimalIds, userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

// handleUpdateEnclosure godoc
//
//	@Id				updateEnclosure
//	@Summary		Update one of the caller's enclosures
//	@Tags			enclosures
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int								true	"Enclosure ID"
//	@Param			enclosure	body	types.UpdateEnclosureV2Payload	true	"Updated enclosure fields"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/enclosures/{id} [put]
func (h *Handler) handleUpdateEnclosure(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	var payload types.UpdateEnclosureV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.UpdateEnclosure(types.Enclosure{
		EnclosureId:   id,
		EnclosureName: payload.EnclosureName,
		HabitatId:     payload.HabitatId,
		Image:         payload.Image,
		Notes:         payload.Notes,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

// handleDeleteEnclosure godoc
//
//	@Id				deleteEnclosure
//	@Summary		Delete one of the caller's enclosures
//	@Description	By default only the enclosure is removed and its animals are left without one. Pass cascade=tasks to also delete the enclosure's tasks, or cascade=animals,tasks to delete its animals and their tasks as well.
//	@Tags			enclosures
//	@Produce		json
//	@Param			id		path	int		true	"Enclosure ID"
//	@Param			cascade	query	string	false	"What else to delete"	Enums(tasks, animals,tasks)
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/enclosures/{id} [delete]
func (h *Handler) handleDeleteEnclosure(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	cascade, err := parseCascade(r.URL.Query().Get("cascade"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	switch cascade {
	case cascadeNone:
		err = h.store.DeleteEnclosureById(id)
	case cascadeTasks:
		err = h.store.DeleteEnclosureAndTasksById(id)
	case cascadeAnimalsAndTasks:
		err = h.store.DeleteEnclosureAndAnimalsAndTasksById(id)
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

type cascadeMode int

const (
	cascadeNone cascadeMode = iota
	cascadeTasks
	cascadeAnimalsAndTasks
)

// parseCascade reads the ?cascade= parameter that replaces v1's separate
// /withtasks and /withanimalsandtasks delete routes. Order is not significant,
// so "tasks,animals" and "animals,tasks" mean the same thing.
func parseCascade(raw string) (cascadeMode, error) {
	if raw == "" {
		return cascadeNone, nil
	}

	var wantAnimals, wantTasks bool

	for part := range strings.SplitSeq(raw, ",") {
		switch strings.TrimSpace(part) {
		case "animals":
			wantAnimals = true
		case "tasks":
			wantTasks = true
		case "":
			continue
		default:
			return cascadeNone, fmt.Errorf("invalid cascade value %q: expected tasks or animals,tasks", raw)
		}
	}

	switch {
	case wantAnimals && wantTasks:
		return cascadeAnimalsAndTasks, nil
	case wantTasks:
		return cascadeTasks, nil
	default:
		// Deleting animals while keeping their tasks would orphan those tasks,
		// and no store method supports it.
		return cascadeNone, fmt.Errorf("invalid cascade value %q: animals cannot be deleted without tasks", raw)
	}
}
