package habitat

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
	store     types.HabitatStore
	userStore types.UserStore
}

func NewHandler(store types.HabitatStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// public routes
	router.HandleFunc("/habitat", h.handleGetHabitats).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/admin/habitat", auth.WithJWTAuth(h.handleAdminCreateHabitat, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/habitat", auth.WithJWTAuth(h.handleAdminUpdateHabitat, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/habitat", auth.WithJWTAuth(h.handleAdminDeleteHabitatById, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetHabitats(w http.ResponseWriter, r *http.Request) {
	habitatsList, err := h.store.GetHabitats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, habitatsList)

}

func (h *Handler) handleAdminCreateHabitat(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var habitat types.CreateHabitatPayload
	if err := utils.ParseJSON(r, &habitat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(habitat); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if habitat exists
	_, err := h.store.GetHabitatByName(habitat.HabitatName)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("habitat with name %s already exists", habitat.HabitatName))
		return
	}

	// if it doesn't exist, create new habitat
	err = h.store.CreateHabitat(types.Habitat{
		HabitatName:    habitat.HabitatName,
		HabitatDesc:    habitat.HabitatDesc,
		Image:          habitat.Image,
		Humidity:       habitat.Humidity,
		DayTempRange:   habitat.DayTempRange,
		NightTempRange: habitat.NightTempRange,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleAdminUpdateHabitat(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var habitat types.UpdateHabitatPayload
	if err := utils.ParseJSON(r, &habitat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(habitat); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// update habitat
	err := h.store.UpdateHabitat(types.Habitat(habitat))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteHabitatById(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var deleteHabitatPayload types.HabitatIdPayload
	if err := utils.ParseJSON(r, &deleteHabitatPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(deleteHabitatPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete habitat
	err := h.store.DeleteHabitatById(deleteHabitatPayload.HabitatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
