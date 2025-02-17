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
	router.HandleFunc("/enclosure", h.handleGetEnclosures).Methods(http.MethodGet)
	router.HandleFunc("/enclosure/family", auth.WithJWTAuth(h.handleGetUserEnclosures, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/enclosure", h.handleCreateEnclosure).Methods(http.MethodPost)
	router.HandleFunc("/animal/family", auth.WithJWTAuth(h.handleCreateUserEnclosure, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleGetEnclosures(w http.ResponseWriter, r *http.Request) {
	enclosureList, err := h.store.GetEnclosures()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleGetUserEnclosures(w http.ResponseWriter, r *http.Request) { //TODO
	// get userId
	userID := auth.GetuserIdFromContext(r.Context()) //start here

	enclosureList, err := h.store.GetEnclosures()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
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

func (h *Handler) handleCreateUserEnclosure(w http.ResponseWriter, r *http.Request) { //TODO
	// get userId
	userID := auth.GetuserIdFromContext(r.Context()) //start here

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
