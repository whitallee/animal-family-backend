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
	// user routes
	router.HandleFunc("/enclosure", auth.WithJWTAuth(h.handleUserGetEnclosures, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/enclosure/id", auth.WithJWTAuth(h.handleUserGetEnclosureById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/enclosure", auth.WithJWTAuth(h.handleUserCreateEnclosure, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/enclosure/withanimals", auth.WithJWTAuth(h.handleUserCreateEnclosureWithAnimals, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/enclosure", auth.WithJWTAuth(h.handleUserUpdateEnclosure, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/enclosure/id", auth.WithJWTAuth(h.handleUserDeleteEnclosureById, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/enclosure/id/withanimals", auth.WithJWTAuth(h.handleUserDeleteEnclosureWithAnimalsById, h.userStore)).Methods(http.MethodDelete)

	// admin routes
	router.HandleFunc("/admin/enclosure", auth.WithJWTAuth(h.handleAdminGetEnclosures, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/enclosure/id", auth.WithJWTAuth(h.handleAdminGetEnclosureById, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/enclosure/byuser", auth.WithJWTAuth(h.handleAdminGetEnclosuresByUser, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/admin/enclosure", auth.WithJWTAuth(h.handleAdminCreateEnclosure, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/enclosure/withanimals", auth.WithJWTAuth(h.handleAdminCreateEnclosureWithAnimals, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/admin/enclosure", auth.WithJWTAuth(h.handleAdminUpdateEnclosure, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/enclosure/owner", auth.WithJWTAuth(h.handleAdminUpdateEnclosureOwner, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/admin/enclosure/id", auth.WithJWTAuth(h.handleAdminDeleteEnclosureById, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/admin/enclosure/id/withanimals", auth.WithJWTAuth(h.handleAdminDeleteEnclosureWithAnimalsById, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleAdminGetEnclosures(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get enclosures
	enclosureList, err := h.store.GetEnclosures()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleAdminGetEnclosuresByUser(w http.ResponseWriter, r *http.Request) {
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

	// get enclosures
	enclosureList, err := h.store.GetEnclosuresByUserId(userIdPayload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleUserGetEnclosures(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get enclosures by user
	enclosureList, err := h.store.GetEnclosuresByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)
}

func (h *Handler) handleAdminGetEnclosureById(w http.ResponseWriter, r *http.Request) {
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

	// get enclosure
	enclosure, err := h.store.GetEnclosureById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosure)
}

func (h *Handler) handleUserGetEnclosureById(w http.ResponseWriter, r *http.Request) {
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
	_, err := h.store.GetEnclosureUserByIds(enclosureIdPayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, get enclosure
	enclosure, err := h.store.GetEnclosureById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosure)
}

func (h *Handler) handleAdminCreateEnclosure(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var enclosure types.CreateEnclosureWithOwnerPayload
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

	// check if enclosure exists
	_, err := h.store.GetEnclosureByNameAndHabitatWithUserId(enclosure.EnclosureName, enclosure.HabitatId, enclosure.UserID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s and habitat id %d already exists", enclosure.EnclosureName, enclosure.HabitatId))
		return
	}

	// if it doesn't exist, create new enclosure
	err = h.store.CreateEnclosure(types.Enclosure{
		EnclosureName: enclosure.EnclosureName,
		Image:         enclosure.Image,
		Notes:         enclosure.Notes,
		HabitatId:     enclosure.HabitatId,
	}, enclosure.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUserCreateEnclosure(w http.ResponseWriter, r *http.Request) {
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

	// check if enclosure exists
	_, err := h.store.GetEnclosureByNameAndHabitatWithUserId(enclosure.EnclosureName, enclosure.HabitatId, userID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s and habitat id %d already exists", enclosure.EnclosureName, enclosure.HabitatId))
		return
	}

	// if it doesn't exist, create new enclosure with userID
	err = h.store.CreateEnclosure(types.Enclosure{
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

func (h *Handler) handleAdminCreateEnclosureWithAnimals(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var enclosureWithAnimals types.CreateEnclosureWithOwnerWithAnimalsPayload
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

	// check if enclosure exists
	_, err := h.store.GetEnclosureByNameAndHabitatWithUserId(enclosureWithAnimals.EnclosureName, enclosureWithAnimals.HabitatId, enclosureWithAnimals.UserID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s and habitat id %d already exists", enclosureWithAnimals.EnclosureName, enclosureWithAnimals.HabitatId))
		return
	}

	// if it doesn't exist, create new enclosure with animals
	err = h.store.CreateEnclosureWithAnimals(types.Enclosure{
		EnclosureName: enclosureWithAnimals.EnclosureName,
		Image:         enclosureWithAnimals.Image,
		Notes:         enclosureWithAnimals.Notes,
		HabitatId:     enclosureWithAnimals.HabitatId,
	}, enclosureWithAnimals.AnimalIds, enclosureWithAnimals.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUserCreateEnclosureWithAnimals(w http.ResponseWriter, r *http.Request) {
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

	// check if enclosure exists
	_, err := h.store.GetEnclosureByNameAndHabitatWithUserId(enclosureWithAnimals.EnclosureName, enclosureWithAnimals.HabitatId, userID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s and habitat id %d already exists", enclosureWithAnimals.EnclosureName, enclosureWithAnimals.HabitatId))
		return
	}

	// if it doesn't exist, create new enclosure with userID
	err = h.store.CreateEnclosureWithAnimals(types.Enclosure{
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

func (h *Handler) handleAdminUpdateEnclosure(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var enclosure types.UpdateEnclosurePayload
	if err := utils.ParseJSON(r, &enclosure); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosure); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// update enclosure
	err := h.store.UpdateEnclosure(types.Enclosure(enclosure))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminUpdateEnclosureOwner(w http.ResponseWriter, r *http.Request) {
	// get userId and check if admin
	userID := auth.GetuserIdFromContext(r.Context())
	if !auth.IsAdmin(userID) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthoized to access this endpoint"))
	}

	// get JSON payload
	var updateEnclosureOwnerPayload types.UpdateEnclosureOwnerPayload
	if err := utils.ParseJSON(r, &updateEnclosureOwnerPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(updateEnclosureOwnerPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// update enclosureUser
	err := h.store.UpdateEnclosureOwnerWithAnimals(types.EnclosureUser{
		EnclosureId: updateEnclosureOwnerPayload.EnclosureId,
		UserID:      updateEnclosureOwnerPayload.OldUserId}, updateEnclosureOwnerPayload.NewUserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserUpdateEnclosure(w http.ResponseWriter, r *http.Request) {
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var enclosure types.UpdateEnclosurePayload
	if err := utils.ParseJSON(r, &enclosure); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosure); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check for ownership
	_, err := h.store.GetEnclosureUserByIds(enclosure.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// check if enclosure exists under user
	_, err = h.store.GetEnclosureByNameAndHabitatWithUserId(enclosure.EnclosureName, enclosure.HabitatId, userID)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("enclosure with name %s and habitat id %d already exists", enclosure.EnclosureName, enclosure.HabitatId))
		return
	}

	// if ownership exists and no duplicates exist, update enclosure
	err = h.store.UpdateEnclosure(types.Enclosure(enclosure))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteEnclosureById(w http.ResponseWriter, r *http.Request) {
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

	// delete enclosure
	err := h.store.DeleteEnclosureById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserDeleteEnclosureById(w http.ResponseWriter, r *http.Request) {
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
	_, err := h.store.GetEnclosureUserByIds(enclosureIdPayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, delete enclosure
	err = h.store.DeleteEnclosureById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteEnclosureWithAnimalsById(w http.ResponseWriter, r *http.Request) {
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

	// delete enclosure with animals
	err := h.store.DeleteEnclosureAndAnimalsById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleUserDeleteEnclosureWithAnimalsById(w http.ResponseWriter, r *http.Request) {
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
	_, err := h.store.GetEnclosureUserByIds(enclosureIdPayload.EnclosureId, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error checking ownership: %v", err))
		return
	}

	// if ownership exists, delete enclosure with animals
	err = h.store.DeleteEnclosureAndAnimalsById(enclosureIdPayload.EnclosureId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
