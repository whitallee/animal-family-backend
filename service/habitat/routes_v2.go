package habitat

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// RegisterV2Routes mounts the v2 habitat routes.
//
// Differences from v1: plural collection path, IDs as path parameters instead
// of JSON body fields, and admin gating via auth.RequireAdmin rather than an
// inline check in each handler.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	router.HandleFunc("/habitats", h.handleListHabitats).Methods(http.MethodGet)

	router.HandleFunc("/habitats", auth.WithJWTAuth(auth.RequireAdmin(h.handleCreateHabitat), h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/habitats/{id}", auth.WithJWTAuth(auth.RequireAdmin(h.handleUpdateHabitat), h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/habitats/{id}", auth.WithJWTAuth(auth.RequireAdmin(h.handleDeleteHabitat), h.userStore)).Methods(http.MethodDelete)
}

// handleListHabitats godoc
//
//	@Id				listHabitats
//	@Summary		List all habitats
//	@Description	Returns every habitat. Habitats are global reference data, not user-owned, so this endpoint is public.
//	@Tags			habitats
//	@Produce		json
//	@Success		200	{array}		types.Habitat
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/habitats [get]
func (h *Handler) handleListHabitats(w http.ResponseWriter, r *http.Request) {
	habitatsList, err := h.store.GetHabitats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, habitatsList)
}

// handleCreateHabitat godoc
//
//	@Id				createHabitat
//	@Summary		Create a habitat
//	@Description	Admin only. Fails if a habitat with the same name already exists.
//	@Tags			habitats
//	@Accept			json
//	@Produce		json
//	@Param			habitat	body	types.CreateHabitatPayload	true	"Habitat to create"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/habitats [post]
func (h *Handler) handleCreateHabitat(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateHabitatPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if _, err := h.store.GetHabitatByName(payload.HabitatName); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("habitat with name %s already exists", payload.HabitatName))
		return
	}

	err := h.store.CreateHabitat(types.Habitat{
		HabitatName:    payload.HabitatName,
		HabitatDesc:    payload.HabitatDesc,
		Image:          payload.Image,
		Humidity:       payload.Humidity,
		DayTempRange:   payload.DayTempRange,
		NightTempRange: payload.NightTempRange,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

// handleUpdateHabitat godoc
//
//	@Id				updateHabitat
//	@Summary		Update a habitat
//	@Description	Admin only. Replaces all fields of the habitat identified by the path parameter.
//	@Tags			habitats
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Habitat ID"
//	@Param			habitat	body	types.UpdateHabitatV2Payload	true	"Updated habitat fields"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		404	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/habitats/{id} [put]
func (h *Handler) handleUpdateHabitat(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDParam(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.UpdateHabitatV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if _, err := h.store.GetHabitatById(id); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("habitat with id %d not found", id))
		return
	}

	err = h.store.UpdateHabitat(types.Habitat{
		HabitatId:      id,
		HabitatName:    payload.HabitatName,
		HabitatDesc:    payload.HabitatDesc,
		Image:          payload.Image,
		Humidity:       payload.Humidity,
		DayTempRange:   payload.DayTempRange,
		NightTempRange: payload.NightTempRange,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

// handleDeleteHabitat godoc
//
//	@Id				deleteHabitat
//	@Summary		Delete a habitat
//	@Description	Admin only.
//	@Tags			habitats
//	@Produce		json
//	@Param			id	path	int	true	"Habitat ID"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		404	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/habitats/{id} [delete]
func (h *Handler) handleDeleteHabitat(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDParam(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := h.store.GetHabitatById(id); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("habitat with id %d not found", id))
		return
	}

	if err := h.store.DeleteHabitatById(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
