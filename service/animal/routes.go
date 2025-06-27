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
	store          types.AnimalStore
	userStore      types.UserStore
	enclosureStore types.EnclosureStore
}

func NewHandler(store types.AnimalStore, userStore types.UserStore, enclosureStore types.EnclosureStore) *Handler {
	return &Handler{store: store, userStore: userStore, enclosureStore: enclosureStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserGetAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal/byid", auth.WithJWTAuth(h.handleUserGetAnimalById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal/byenclosure", auth.WithJWTAuth(h.handleUserGetAnimalsByEnclosure, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserUpdateAnimal, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserDeleteAnimal, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminGetAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byid", auth.WithJWTAuth(h.handleAdminGetAnimalById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byenclosure", auth.WithJWTAuth(h.handleAdminGetAnimalsByEnclosure, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byuser", auth.WithJWTAuth(h.handleAdminGetAnimalsByUser, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminUpdateAnimal, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/animal/owner", auth.WithJWTAuth(h.handleAdminUpdateAnimalOwner, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminDeleteAnimal, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleAdminCreateAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var animal types.CreateAnimalWithOwnerPayload
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
	_, err := h.store.GetAnimalByNameAndSpeciesWithUserId(animal.AnimalName, animal.SpeciesId, animal.UserID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal with name %s and species id %d already exists", animal.AnimalName, animal.SpeciesId))
		return
	}

	// if it doesn't exist, create new animal
	err = h.store.CreateAnimal(types.Animal{
		AnimalName:      animal.AnimalName,
		SpeciesId:       animal.SpeciesId,
		EnclosureId:     animal.EnclosureId,
		Image:           animal.Image,
		Gender:          animal.Gender,
		Dob:             animal.Dob,
		PersonalityDesc: animal.PersonalityDesc,
		DietDesc:        animal.DietDesc,
		RoutineDesc:     animal.RoutineDesc,
		ExtraNotes:      animal.ExtraNotes,
	}, animal.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUserCreateAnimal(w http.ResponseWriter, r *http.Request) {
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

	// check if animal exists
	_, err := h.store.GetAnimalByNameAndSpeciesWithUserId(animal.AnimalName, animal.SpeciesId, userID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal with name %s and species id %d already exists", animal.AnimalName, animal.SpeciesId))
		return
	}

	// if it doesn't exist, create new animal
	err = h.store.CreateAnimal(types.Animal{
		AnimalName:      animal.AnimalName,
		SpeciesId:       animal.SpeciesId,
		EnclosureId:     animal.EnclosureId,
		Image:           animal.Image,
		Gender:          animal.Gender,
		Dob:             animal.Dob,
		PersonalityDesc: animal.PersonalityDesc,
		DietDesc:        animal.DietDesc,
		RoutineDesc:     animal.RoutineDesc,
		ExtraNotes:      animal.ExtraNotes,
	}, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleAdminUpdateAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var animal types.UpdateAnimalPayload
	if err := utils.ParseJSON(r, &animal); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animal); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// get owner's userID
	animalUser, err := h.store.GetAnimalUserByAnimalId(animal.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if dupe of animal exists under user
	_, err = h.store.GetAnimalByNameAndSpeciesWithUserId(animal.AnimalName, animal.SpeciesId, animalUser.UserID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal with name %s and species id %d already exists", animal.AnimalName, animal.SpeciesId))
		return
	}

	// update animal if no dupe exists
	err = h.store.UpdateAnimal(types.Animal(animal))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserUpdateAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var animal types.UpdateAnimalPayload
	if err := utils.ParseJSON(r, &animal); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animal); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check for ownership
	_, err := h.store.GetAnimalUserByIds(animal.AnimalId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// check if dupe of animal exists under user
	_, err = h.store.GetAnimalByNameAndSpeciesWithUserId(animal.AnimalName, animal.SpeciesId, userID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal with name %s and species id %d already exists", animal.AnimalName, animal.SpeciesId))
		return
	}

	// if ownership exists, update animal
	err = h.store.UpdateAnimal(types.Animal(animal))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminUpdateAnimalOwner(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var updateAnimalOwnerPayload types.UpdateAnimalOwnerPayload
	if err := utils.ParseJSON(r, &updateAnimalOwnerPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(updateAnimalOwnerPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// update animalUser
	err := h.store.UpdateAnimalOwner(types.AnimalUser{
		AnimalId: updateAnimalOwnerPayload.AnimalId,
		UserID:   updateAnimalOwnerPayload.OldUserId}, updateAnimalOwnerPayload.NewUserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminGetAnimals(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized to access this endpoint"))
	}

	// get animals
	animalList, err := h.store.GetAnimals()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleAdminGetAnimalsByUser(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var userIdPayload types.UserIDPayload
	if err := utils.ParseJSON(r, &userIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(userIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get animals
	animalList, err := h.store.GetAnimalsByUserId(userIdPayload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleUserGetAnimals(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get animals
	animalList, err := h.store.GetAnimalsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleAdminGetAnimalById(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var animalIdPayload types.AnimalIdPayload
	if err := utils.ParseJSON(r, &animalIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animalIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get animal
	animal, err := h.store.GetAnimalById(animalIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animal)
}

func (h *Handler) handleUserGetAnimalById(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var animalIdPayload types.AnimalIdPayload
	if err := utils.ParseJSON(r, &animalIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animalIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check for ownership
	_, err := h.store.GetAnimalUserByIds(animalIdPayload.AnimalId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, get animal
	animal, err := h.store.GetAnimalById(animalIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animal)
}

func (h *Handler) handleAdminGetAnimalsByEnclosure(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var enclosureIdPayload types.EnclosureIdPayload
	if err := utils.ParseJSON(r, &enclosureIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosureIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get animals
	animalList, err := h.store.GetAnimalsByEnclosureId(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)

}

func (h *Handler) handleUserGetAnimalsByEnclosure(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosureIdPayload types.EnclosureIdPayload
	if err := utils.ParseJSON(r, &enclosureIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosureIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check for ownership
	_, err := h.enclosureStore.GetEnclosureUserByIds(enclosureIdPayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, get animals
	animalList, err := h.store.GetAnimalsByEnclosureId(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, animalList)
}

func (h *Handler) handleAdminDeleteAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var animaIdPayload types.AnimalIdPayload
	if err := utils.ParseJSON(r, &animaIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animaIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete animal
	err := h.store.DeleteAnimalById(animaIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserDeleteAnimal(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var animaIdPayload types.AnimalIdPayload
	if err := utils.ParseJSON(r, &animaIdPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(animaIdPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check for ownership
	_, err := h.store.GetAnimalUserByIds(animaIdPayload.AnimalId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists delete animal
	err = h.store.DeleteAnimalById(animaIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
