package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/service/auth"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleCreateUser).Methods(http.MethodPost)
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/delete-user", auth.WithJWTAuth(h.handleSelfDeleteUserById, h.store)).Methods(http.MethodDelete)
	router.HandleFunc("/delete-user/id", auth.WithJWTAuth(h.handleAdminDeleteUserById, h.store)).Methods(http.MethodDelete)
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// hash password before user creation
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesn't exist, create new user
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil { // I changed &user to user now back to &user
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// find user
	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	// compare password hash
	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleSelfDeleteUserById(w http.ResponseWriter, r *http.Request) { //untested
	// get userId
	userID := auth.GetuserIdFromContext(r.Context())

	// delete user
	err := h.store.DeleteUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error, please try again, user may not be deleted"))
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) handleAdminDeleteUserById(w http.ResponseWriter, r *http.Request) { //untested
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

	// delete user
	err := h.store.DeleteUserById(userIdPayload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error, please try again, user may not be deleted"))
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
