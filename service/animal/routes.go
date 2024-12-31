package animal

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.AnimalStore
}

func NewHandler(store types.AnimalStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/animal", h.handleGetAnimals).Methods(http.MethodGet)
	router.HandleFunc("/animal", h.handleCreateAnimal).Methods(http.MethodPost)
}

func (h *Handler) handleGetAnimals(w http.ResponseWriter, r *http.Request) {
	animalList, err := h.store.GetAnimals()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}
func (h *Handler) handleCreateAnimal(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var animal types.CreateAnimalPayload
	if err := utils.ParseJSON(r, &animal); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animal); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if animal exists
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
	err := h.store.CreateAnimal(types.Animal{
		AnimalName:  animal.AnimalName,
		SpeciesId:   animal.SpeciesId,
		EnclosureId: animal.EnclosureId,
		Image:       animal.Image,
		Notes:       animal.Notes,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
