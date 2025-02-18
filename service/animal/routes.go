package animal

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// Handler struct contains all of the stores needed for the service
type Handler struct {
	store     types.AnimalStore
	userStore types.UserStore
}

func NewHandler(store types.AnimalStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/animal", h.handleCreateAnimal).Methods(http.MethodPost)
	router.HandleFunc("/animal/family", auth.WithJWTAuth(h.handleCreateAnimalByUserId, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/animal", h.handleGetAnimals).Methods(http.MethodGet)
	router.HandleFunc("/animal/family", auth.WithJWTAuth(h.handleGetAnimalsByUserId, h.userStore)).Methods(http.MethodGet)
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

	// TODO check if animal exists
	// _, err := h.store.GetAnimalByName(species.comName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with common name %s already exists", species.comName))
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

func (h *Handler) handleCreateAnimalByUserId(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

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

	// TODO check if animal exists under the userID with same name
	// _, err := h.store.GetAnimalByNameAndUserID(animal.name, userID)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal named %s already exists in user animal family", animal.name))
	// 	return
	// }

	// if it doesn't exist, create new animal
	err := h.store.CreateAnimalByUserId(types.Animal{
		AnimalName:  animal.AnimalName,
		SpeciesId:   animal.SpeciesId,
		EnclosureId: animal.EnclosureId,
		Image:       animal.Image,
		Notes:       animal.Notes,
	}, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetAnimals(w http.ResponseWriter, r *http.Request) {
	animalList, err := h.store.GetAnimals()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleGetAnimalsByUserId(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	animalList, err := h.store.GetAnimalsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}
