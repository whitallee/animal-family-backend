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

type Handler struct {
	store     types.SpeciesStore
	userStore types.UserStore
}

func NewHandler(store types.SpeciesStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/species", h.handleGetSpecies).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/admin/species", auth.WithJWTAuth(h.handleAdminCreateSpecies, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/species", auth.WithJWTAuth(h.handleAdminUpdateSpecies, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/species", auth.WithJWTAuth(h.handleAdminDeleteSpeciesById, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetSpecies(w http.ResponseWriter, r *http.Request) {
	speciesList, err := h.store.GetSpecies()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, speciesList)

}

func (h *Handler) handleAdminCreateSpecies(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var species types.CreateSpeciesPayload
	if err := utils.ParseJSON(r, &species); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(species); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if species exists already
	_, err := h.store.GetSpeciesByComName(species.ComName)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with common name %s already exists", species.ComName))
		return
	}
	_, err = h.store.GetSpeciesBySciName(species.SciName)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with scientific name %s already exists", species.SciName))
		return
	}

	// if it doesn't exist, create new species
	err = h.store.CreateSpecies(types.Species{
		ComName:     species.ComName,
		SciName:     species.SciName,
		SpeciesDesc: species.SpeciesDesc,
		Image:       species.Image,
		HabitatId:   species.HabitatId,
		BaskTemp:    species.BaskTemp,
		Diet:        species.Diet,
		Sociality:   species.Sociality,
		ExtraCare:   species.ExtraCare,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleAdminUpdateSpecies(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var species types.UpdateSpeciesPayload
	if err := utils.ParseJSON(r, &species); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(species); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update species
	err := h.store.CreateSpecies(types.Species(species))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteSpeciesById(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var deleteSpeciesPayload types.SpeciesIdPayload
	if err := utils.ParseJSON(r, &deleteSpeciesPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(deleteSpeciesPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete species
	err := h.store.DeleteSpeciesById(deleteSpeciesPayload.SpeciesId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
