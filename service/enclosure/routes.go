package enclosure

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
	store     types.EnclosureStore
	userStore types.UserStore
}

func NewHandler(store types.EnclosureStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/enclosure", h.handleCreateEnclosure).Methods(http.MethodPost)
	router.HandleFunc("/family/enclosure", auth.WithJWTAuth(h.handleCreateEnclosureByUserID, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/family/enclosure/withanimals", auth.WithJWTAuth(h.handleCreateEnclosureWithAnimalsByUserID, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/enclosure", h.handleGetEnclosures).Methods(http.MethodGet)
	router.HandleFunc("/family/enclosure", auth.WithJWTAuth(h.handleGetEnclosuresByUserId, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/family/enclosure/id", auth.WithJWTAuth(h.handleGetEnclosureByIdWithUserId, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/family/enclosure/id", auth.WithJWTAuth(h.handleDeleteEnclosureByIdWithUserId, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/family/enclosure/id/withanimals", auth.WithJWTAuth(h.handleDeleteEnclosureWithAnimalsByIdWithUserId, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleCreateEnclosure(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var enclosure types.CreateEnclosurePayload
	if err := utils.ParseJSON(r, &enclosure); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosure); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// TODO check if enclosure exists

	// if it doesn't exist, create new enclosure
	err := h.store.CreateEnclosure(types.Enclosure{
		EnclosureName: enclosure.EnclosureName,
		Image:         enclosure.Image,
		Notes:         enclosure.Notes,
		HabitatId:     enclosure.HabitatId,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleCreateEnclosureByUserID(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosure types.CreateEnclosurePayload
	if err := utils.ParseJSON(r, &enclosure); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosure); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// TODO check if enclosure exists

	// if it doesn't exist, create new enclosure with userID
	err := h.store.CreateEnclosureByUserId(types.Enclosure{
		EnclosureName: enclosure.EnclosureName,
		Image:         enclosure.Image,
		Notes:         enclosure.Notes,
		HabitatId:     enclosure.HabitatId,
	}, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleCreateEnclosureWithAnimalsByUserID(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosureWithAnimals types.CreateEnclosureWithAnimalsPayload
	if err := utils.ParseJSON(r, &enclosureWithAnimals); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosureWithAnimals); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// TODO check if enclosure exists

	// if it doesn't exist, create new enclosure with userID
	err := h.store.CreateEnclosureWithAnimalsByUserId(types.Enclosure{
		EnclosureName: enclosureWithAnimals.EnclosureName,
		Image:         enclosureWithAnimals.Image,
		Notes:         enclosureWithAnimals.Notes,
		HabitatId:     enclosureWithAnimals.HabitatId,
	}, enclosureWithAnimals.AnimalIds, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetEnclosures(w http.ResponseWriter, r *http.Request) {
	enclosureList, err := h.store.GetEnclosures()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleGetEnclosuresByUserId(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	enclosureList, err := h.store.GetEnclosuresByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleGetEnclosureByIdWithUserId(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosureId types.EnclosureIdPayload
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

	enclosure, err := h.store.GetEnclosureByIdWithUserId(enclosureId.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosure)
}

func (h *Handler) handleDeleteEnclosureByIdWithUserId(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var deleteEnclosurePayload types.EnclosureIdPayload
	if err := utils.ParseJSON(r, &deleteEnclosurePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(deleteEnclosurePayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// TODO check if enclosure exists

	// if it doesn't exist, create new enclosure with userID
	err := h.store.DeleteEnclosureByIdWithUserId(deleteEnclosurePayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleDeleteEnclosureWithAnimalsByIdWithUserId(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var deleteEnclosurePayload types.EnclosureIdPayload
	if err := utils.ParseJSON(r, &deleteEnclosurePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(deleteEnclosurePayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// TODO check if enclosure exists

	// if it doesn't exist, create new enclosure with userID
	err := h.store.DeleteEnclosureAndAnimalsByIdWithUserId(deleteEnclosurePayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
