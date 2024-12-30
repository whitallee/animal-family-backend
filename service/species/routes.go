package species

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.SpeciesStore
}

func NewHandler(store types.SpeciesStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/species", h.handleGetSpecies).Methods(http.MethodGet)
	router.HandleFunc("/species", h.handleCreateSpecies).Methods(http.MethodPost)
}

func (h *Handler) handleGetSpecies(w http.ResponseWriter, r *http.Request) {
	speciesList, err := h.store.GetSpecies()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, speciesList)

}
func (h *Handler) handleCreateSpecies(w http.ResponseWriter, r *http.Request) {
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

	// check if species exists
	// _, err := h.store.GetSpeciesByComName(species.comName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with common name %s already exists", species.comName))
	// 	return
	// }
	// _, err := h.store.GetSpeciesBySciName(species.sciName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with scientific name %s already exists", species.sciName))
	// 	return
	// }

	// if it doesn't exist, create new species
	err := h.store.CreateSpecies(types.Species{
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
