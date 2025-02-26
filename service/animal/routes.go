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
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleAdminCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/family/animal", auth.WithJWTAuth(h.handleCreateAnimalByUserId, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleAdminGetAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/family/animal", auth.WithJWTAuth(h.handleGetAnimalsByUserId, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/family/animal/byenclosure", auth.WithJWTAuth(h.handleGetAnimalsByEnclosureIdWithUserId, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/family/animal", auth.WithJWTAuth(h.handleDeleteAnimalByIdWithUserId, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleAdminCreateAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

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

func (h *Handler) handleAdminGetAnimals(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get animals
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

func (h *Handler) handleGetAnimalsByEnclosureIdWithUserId(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosureId types.GetAnimalsByEnclosureIdWithUserIdPayload
	if err := utils.ParseJSON(r, &enclosureId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosureId); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	animalList, err := h.store.GetAnimalsByEnclosureIdWithUserId(enclosureId.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleDeleteAnimalByIdWithUserId(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var deleteAnimalPayload types.DeleteAnimalByIdWithUserIdPayload
	if err := utils.ParseJSON(r, &deleteAnimalPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(deleteAnimalPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete animal
	err := h.store.DeleteAnimalByIdWithUserId(deleteAnimalPayload.AnimalId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
