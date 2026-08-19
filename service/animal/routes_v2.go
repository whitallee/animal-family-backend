package animal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

// RegisterV2Routes mounts the v2 animal routes.
//
// Differences from v1: plural collection path, IDs as path parameters, the
// separate /byenclosure route folded into GET /animals?enclosureId=, and the
// /withtasks delete route replaced by ?cascade=tasks. Responses use
// types.AnimalResponse rather than types.Animal — see the note on that type.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	owned := func(next http.HandlerFunc) http.HandlerFunc {
		return auth.WithJWTAuth(auth.RequireOwnership("id", h.store.UserOwnsAnimal, next), h.userStore)
	}

	router.HandleFunc("/animals", auth.WithJWTAuth(h.handleListAnimals, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/animals", auth.WithJWTAuth(h.handleCreateAnimal, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/animals/{id}", owned(h.handleGetAnimal)).Methods(http.MethodGet)
	router.HandleFunc("/animals/{id}", owned(h.handleUpdateAnimal)).Methods(http.MethodPut)
	router.HandleFunc("/animals/{id}", owned(h.handleDeleteAnimal)).Methods(http.MethodDelete)

	// Memorial state is a sub-resource rather than fields on the animal, so an
	// ordinary edit cannot clear it by omission.
	router.HandleFunc("/animals/{id}/memorial", owned(h.handleSetAnimalMemorial)).Methods(http.MethodPut)
	router.HandleFunc("/animals/{id}/memorial", owned(h.handleClearAnimalMemorial)).Methods(http.MethodDelete)
}

// handleListAnimals godoc
//
//	@Id				listAnimals
//	@Summary		List the caller's animals
//	@Description	Pass enclosureId to return only the animals in that enclosure. The enclosure must belong to the caller.
//	@Tags			animals
//	@Produce		json
//	@Param			enclosureId	query	int	false	"Only return animals in this enclosure"
//	@Success		200	{array}		types.AnimalResponse
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals [get]
func (h *Handler) handleListAnimals(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	animals, ok := h.listAnimalsFor(w, r, userID)
	if !ok {
		return
	}

	responses, err := types.NewAnimalResponses(animals)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, responses)
}

// listAnimalsFor applies the optional enclosureId filter that replaces v1's
// /animal/byenclosure route. It writes the response itself on failure and
// reports whether the caller should continue, because the failure modes need
// three different statuses.
func (h *Handler) listAnimalsFor(w http.ResponseWriter, r *http.Request, userID int) ([]*types.Animal, bool) {
	raw := r.URL.Query().Get("enclosureId")

	if raw == "" {
		animals, err := h.store.GetAnimalsByUserId(userID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return nil, false
		}

		return animals, true
	}

	enclosureId, err := strconv.Atoi(raw)
	if err != nil || enclosureId < 1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid enclosureId: must be a positive integer"))
		return nil, false
	}

	// Without this check, anyone could enumerate another user's animals by
	// guessing enclosure IDs.
	owned, err := h.enclosureStore.UserOwnsEnclosure(enclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return nil, false
	}
	if !owned {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you do not have access to this enclosure"))
		return nil, false
	}

	animals, err := h.store.GetAnimalsByEnclosureId(enclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return nil, false
	}

	return animals, true
}

// handleGetAnimal godoc
//
//	@Id				getAnimal
//	@Summary		Get one of the caller's animals
//	@Tags			animals
//	@Produce		json
//	@Param			id	path		int	true	"Animal ID"
//	@Success		200	{object}	types.AnimalResponse
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals/{id} [get]
func (h *Handler) handleGetAnimal(w http.ResponseWriter, r *http.Request) {
	// Already parsed and ownership-checked by RequireOwnership.
	id := auth.ResourceIDFromContext(r.Context())

	animal, err := h.store.GetAnimalById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := types.NewAnimalResponse(animal)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// handleCreateAnimal godoc
//
//	@Id				createAnimal
//	@Summary		Create an animal
//	@Tags			animals
//	@Accept			json
//	@Produce		json
//	@Param			animal	body	types.CreateAnimalV2Payload	true	"Animal to create"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals [post]
func (h *Handler) handleCreateAnimal(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.CreateAnimalV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if payload.EnclosureId != nil {
		owned, err := h.enclosureStore.UserOwnsEnclosure(*payload.EnclosureId, userID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !owned {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you do not have access to that enclosure"))
			return
		}
	}

	if _, err := h.store.GetAnimalByNameAndSpeciesWithUserId(payload.AnimalName, payload.SpeciesId, userID); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("animal named %s of that species already exists", payload.AnimalName))
		return
	}

	err := h.store.CreateAnimal(types.Animal{
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
	}, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusCreated)
}

// handleUpdateAnimal godoc
//
//	@Id				updateAnimal
//	@Summary		Update one of the caller's animals
//	@Description	A full replace: every field is set to the value sent, and an omitted field is reset rather than left alone. Memorial state is not affected — use /animals/{id}/memorial for that.
//	@Tags			animals
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Animal ID"
//	@Param			animal	body	types.UpdateAnimalV2Payload	true	"Updated animal fields"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals/{id} [put]
func (h *Handler) handleUpdateAnimal(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.UpdateAnimalV2Payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if payload.EnclosureId != nil {
		owned, err := h.enclosureStore.UserOwnsEnclosure(*payload.EnclosureId, userID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !owned {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you do not have access to that enclosure"))
			return
		}
	}

	// UpdateAnimalDetails leaves the memorial columns alone, so no read of the
	// stored animal is needed to carry them across.
	err := h.store.UpdateAnimalDetails(animalFromUpdate(payload, id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// animalFromUpdate converts the payload into the record to store.
//
// It is a plain field-for-field copy: PUT replaces the animal's details, so
// every value comes from the request and nothing is merged with what was there
// before. Memorial fields are absent by design — see UpdateAnimalV2Payload.
func animalFromUpdate(payload types.UpdateAnimalV2Payload, id int) types.Animal {
	return types.Animal{
		AnimalId:        id,
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
}

// handleSetAnimalMemorial godoc
//
//	@Id				setAnimalMemorial
//	@Summary		Memorialise one of the caller's animals
//	@Description	Records that an animal has died, along with a message, any photos, and the date. Memorialised animals are filtered out of the main views. Sending this again replaces the stored memorial.
//	@Tags			animals
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int								true	"Animal ID"
//	@Param			memorial	body	types.SetAnimalMemorialPayload	true	"Memorial details"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals/{id}/memorial [put]
func (h *Handler) handleSetAnimalMemorial(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	var payload types.SetAnimalMemorialPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	// Stored as a JSON array so it reads back the same way whether or not any
	// photos were supplied; a nil slice would encode as "null".
	photos := payload.MemorialPhotos
	if photos == nil {
		photos = []string{}
	}

	encoded, err := json.Marshal(photos)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not encode memorialPhotos: %w", err))
		return
	}

	err = h.store.SetAnimalMemorial(id, payload.LastMessage, string(encoded), payload.MemorialDate)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// handleClearAnimalMemorial godoc
//
//	@Id				clearAnimalMemorial
//	@Summary		Remove an animal's memorial
//	@Description	Returns a memorialised animal to the living roster, discarding its message and photos.
//	@Tags			animals
//	@Produce		json
//	@Param			id	path	int	true	"Animal ID"
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals/{id}/memorial [delete]
func (h *Handler) handleClearAnimalMemorial(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	if err := h.store.ClearAnimalMemorial(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}

// handleDeleteAnimal godoc
//
//	@Id				deleteAnimal
//	@Summary		Delete one of the caller's animals
//	@Description	Pass cascade=tasks to also delete the animal's tasks, replacing v1's separate /animal/withtasks route.
//	@Tags			animals
//	@Produce		json
//	@Param			id		path	int		true	"Animal ID"
//	@Param			cascade	query	string	false	"What else to delete"	Enums(tasks)
//	@Success		204
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/animals/{id} [delete]
func (h *Handler) handleDeleteAnimal(w http.ResponseWriter, r *http.Request) {
	id := auth.ResourceIDFromContext(r.Context())

	cascade := r.URL.Query().Get("cascade")

	var err error
	switch cascade {
	case "":
		err = h.store.DeleteAnimalById(id)
	case "tasks":
		err = h.store.DeleteAnimalAndTasksById(id)
	default:
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid cascade value %q: expected tasks", cascade))
		return
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}
