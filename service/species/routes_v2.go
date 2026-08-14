package species

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// RegisterV2Routes mounts the v2 species routes.
//
// "species" is both singular and plural, so the collection path is unchanged.
// Differences from v1: IDs as path parameters instead of JSON body fields, and
// admin gating via auth.RequireAdmin rather than an inline check per handler.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	router.HandleFunc("/species", h.handleListSpecies).Methods(http.MethodGet)

	router.HandleFunc("/species", auth.WithJWTAuth(auth.RequireAdmin(h.handleCreateSpecies), h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/species/generate", auth.WithJWTAuth(auth.RequireAdmin(h.handleGenerateSpecies), h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/species/{id}", auth.WithJWTAuth(auth.RequireAdmin(h.handleUpdateSpecies), h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/species/{id}", auth.WithJWTAuth(auth.RequireAdmin(h.handleDeleteSpecies), h.userStore)).Methods(http.MethodDelete)
}

// handleListSpecies godoc
//
//	@Id				listSpecies
//	@Summary		List all species
//	@Description	Returns every species. Species are global reference data, not user-owned, so this endpoint is public.
//	@Tags			species
//	@Produce		json
//	@Success		200	{array}		types.Species
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/species [get]
func (h *Handler) handleListSpecies(w http.ResponseWriter, r *http.Request) {
	speciesList, err := h.store.GetSpecies()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, speciesList)
}

// handleCreateSpecies godoc
//
//	@Id				createSpecies
//	@Summary		Create a species
//	@Description	Admin only. Fails if a species with the same common or scientific name already exists.
//	@Tags			species
//	@Accept			json
//	@Produce		json
//	@Param			species	body	types.CreateSpeciesPayload	true	"Species to create"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/species [post]
func (h *Handler) handleCreateSpecies(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateSpeciesPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if _, err := h.store.GetSpeciesByComName(payload.ComName); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with common name %s already exists", payload.ComName))
		return
	}
	if _, err := h.store.GetSpeciesBySciName(payload.SciName); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with scientific name %s already exists", payload.SciName))
		return
	}

	// Unlike v1's handleAdminCreateSpecies, this persists every validated field.
	// That handler omitted Lifespan, Size, Weight and ConservationStatus, so
	// those were stored empty despite being validate:"required".
	err := h.store.CreateSpecies(types.Species{
		ComName:            payload.ComName,
		SciName:            payload.SciName,
		SpeciesDesc:        payload.SpeciesDesc,
		Image:              payload.Image,
		HabitatId:          payload.HabitatId,
		BaskTemp:           payload.BaskTemp,
		Diet:               payload.Diet,
		Sociality:          payload.Sociality,
		Lifespan:           payload.Lifespan,
		Size:               payload.Size,
		Weight:             payload.Weight,
		ConservationStatus: payload.ConservationStatus,
		ExtraCare:          payload.ExtraCare,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusCreated)
}

// handleGenerateSpecies godoc
//
//	@Id				generateSpecies
//	@Summary		Generate a species from a name
//	@Description	Admin only. Uses an LLM to populate species details and generates an image, uploading it to S3. Fails if a species with that name already exists (case insensitive).
//	@Tags			species
//	@Accept			json
//	@Produce		json
//	@Param			species	body		types.GenerateSpeciesPayload	true	"Common name of the species to generate"
//	@Success		201		{object}	types.Species
//	@Failure		400		{object}	types.ErrorResponse
//	@Failure		403		{object}	types.ErrorResponse
//	@Failure		500		{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/species/generate [post]
func (h *Handler) handleGenerateSpecies(w http.ResponseWriter, r *http.Request) {
	var payload types.GenerateSpeciesPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if _, err := h.store.GetSpeciesByNameCaseInsensitive(payload.Name); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with name '%s' already exists", payload.Name))
		return
	}

	species, err := h.store.GenerateSpeciesFromName(payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, species)
}

// handleUpdateSpecies godoc
//
//	@Id				updateSpecies
//	@Summary		Update a species
//	@Description	Admin only. Replaces all fields of the species identified by the path parameter.
//	@Tags			species
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Species ID"
//	@Param			species	body	types.UpdateSpeciesV2Payload	true	"Updated species fields"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		404	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/species/{id} [put]
func (h *Handler) handleUpdateSpecies(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDParam(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.UpdateSpeciesV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if _, err := h.store.GetSpeciesById(id); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("species with id %d not found", id))
		return
	}

	err = h.store.UpdateSpecies(types.Species{
		SpeciesID:          id,
		ComName:            payload.ComName,
		SciName:            payload.SciName,
		SpeciesDesc:        payload.SpeciesDesc,
		Image:              payload.Image,
		HabitatId:          payload.HabitatId,
		BaskTemp:           payload.BaskTemp,
		Diet:               payload.Diet,
		Sociality:          payload.Sociality,
		Lifespan:           payload.Lifespan,
		Size:               payload.Size,
		Weight:             payload.Weight,
		ConservationStatus: payload.ConservationStatus,
		ExtraCare:          payload.ExtraCare,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// handleDeleteSpecies godoc
//
//	@Id				deleteSpecies
//	@Summary		Delete a species
//	@Description	Admin only.
//	@Tags			species
//	@Produce		json
//	@Param			id	path	int	true	"Species ID"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		404	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/species/{id} [delete]
func (h *Handler) handleDeleteSpecies(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDParam(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := h.store.GetSpeciesById(id); err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("species with id %d not found", id))
		return
	}

	if err := h.store.DeleteSpeciesById(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}
