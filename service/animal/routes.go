package animal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func convertUpdatePayloadToAnimal(payload types.UpdateAnimalPayload, existingAnimal *types.Animal) (types.Animal, error) {
	animal := types.Animal{
		AnimalId:        payload.AnimalId,
		AnimalName:      payload.AnimalName,
		SpeciesId:       payload.SpeciesId,
		EnclosureId:     payload.EnclosureId,
		Image:           payload.Image,
		Gender:          payload.Gender,
		Dob:             payload.Dob,
		PersonalityDesc: payload.PersonalityDesc,
		DietDesc:        payload.DietDesc,
		RoutineDesc:     payload.RoutineDesc,
		ExtraNotes:      payload.ExtraNotes,
	}

	// Handle IsMemorialized - use provided value or preserve existing
	if payload.IsMemorialized != nil {
		animal.IsMemorialized = *payload.IsMemorialized
	} else if existingAnimal != nil {
		animal.IsMemorialized = existingAnimal.IsMemorialized
	} else {
		animal.IsMemorialized = false
	}

	// Handle LastMessage - convert *string to sql.NullString or preserve existing
	if payload.LastMessage != nil {
		animal.LastMessage = sql.NullString{String: *payload.LastMessage, Valid: true}
	} else if existingAnimal != nil {
		animal.LastMessage = existingAnimal.LastMessage
	} else {
		animal.LastMessage = sql.NullString{Valid: false}
	}

	// Handle MemorialPhotos - marshal []string to JSON or preserve existing
	if payload.MemorialPhotos != nil {
		memorialPhotosJSON, err := json.Marshal(payload.MemorialPhotos)
		if err != nil {
			return animal, fmt.Errorf("failed to marshal memorialPhotos: %w", err)
		}
		animal.MemorialPhotos = sql.NullString{String: string(memorialPhotosJSON), Valid: true}
	} else if existingAnimal != nil {
		animal.MemorialPhotos = existingAnimal.MemorialPhotos
	} else {
		animal.MemorialPhotos = sql.NullString{Valid: false}
	}

	// Handle MemorialDate - use provided value if not zero, otherwise preserve existing or let DB set default
	if !payload.MemorialDate.IsZero() {
		animal.MemorialDate = payload.MemorialDate
	} else if existingAnimal != nil {
		animal.MemorialDate = existingAnimal.MemorialDate
	} else {
		// Database will set default (CURRENT_DATE) for new records
		animal.MemorialDate = time.Time{}
	}

	return animal, nil
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserGetAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal/byid", auth.WithJWTAuth(h.handleUserGetAnimalById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal/byenclosure", auth.WithJWTAuth(h.handleUserGetAnimalsByEnclosure, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserUpdateAnimal, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/animal", auth.WithJWTAuth(h.handleUserDeleteAnimal, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/animal/withtasks", auth.WithJWTAuth(h.handleUserDeleteAnimalWithTasks, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminGetAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byid", auth.WithJWTAuth(h.handleAdminGetAnimalById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byenclosure", auth.WithJWTAuth(h.handleAdminGetAnimalsByEnclosure, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal/byuser", auth.WithJWTAuth(h.handleAdminGetAnimalsByUser, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminUpdateAnimal, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/animal/owner", auth.WithJWTAuth(h.handleAdminUpdateAnimalOwner, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/animal", auth.WithJWTAuth(h.handleAdminDeleteAnimal, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/admin/animal/withtasks", auth.WithJWTAuth(h.handleAdminDeleteAnimalWithTasks, h.userStore)).Methods(http.MethodDelete)
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
		IsMemorialized:  false,
		LastMessage:     sql.NullString{Valid: false},
		MemorialPhotos:  sql.NullString{Valid: false},
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
		IsMemorialized:  false,
		LastMessage:     sql.NullString{Valid: false},
		MemorialPhotos:  sql.NullString{Valid: false},
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

	// get existing animal to preserve values for optional fields
	existingAnimal, err := h.store.GetAnimalById(animal.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// convert payload to Animal type
	animalToUpdate, err := convertUpdatePayloadToAnimal(animal, existingAnimal)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// update animal if no dupe exists
	err = h.store.UpdateAnimal(animalToUpdate)
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
	existingAnimalByName, err := h.store.GetAnimalByNameAndSpeciesWithUserId(animal.AnimalName, animal.SpeciesId, userID)
	if err == nil && existingAnimalByName.AnimalId != animal.AnimalId {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal with name %s and species id %d already exists", animal.AnimalName, animal.SpeciesId))
		return
	}

	// get existing animal to preserve values for optional fields
	existingAnimal, err := h.store.GetAnimalById(animal.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// convert payload to Animal type
	animalToUpdate, err := convertUpdatePayloadToAnimal(animal, existingAnimal)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// if ownership exists, update animal
	err = h.store.UpdateAnimal(animalToUpdate)
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

func (h *Handler) handleUserDeleteAnimalWithTasks(w http.ResponseWriter, r *http.Request) {
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

	// if ownership exists, delete animal with tasks
	err = h.store.DeleteAnimalAndTasksById(animaIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteAnimalWithTasks(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
		return
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

	// delete animal with tasks
	err := h.store.DeleteAnimalAndTasksById(animaIdPayload.AnimalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
