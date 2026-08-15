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

// RegisterV2Routes mounts the v2 user routes.
//
// The two admin delete routes from v1 are not carried over: they duplicated
// DELETE /users/me with a different way of naming the target. Admin access can
// return later as a permission check on these routes rather than a parallel
// tree.
//
// Responses use types.UserResponse, never types.User — see the note on that
// type for why returning the domain struct produces the wrong wire format.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	router.HandleFunc("/users/register", h.handleRegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/users/login", h.handleLoginUser).Methods(http.MethodPost)
	router.HandleFunc("/users/refresh-token", auth.WithJWTAuth(h.handleRefreshToken, h.store)).Methods(http.MethodPost)
	router.HandleFunc("/users/me", auth.WithJWTAuth(h.handleGetCurrentUser, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/users/me", auth.WithJWTAuth(h.handleDeleteCurrentUser, h.store)).Methods(http.MethodDelete)
}

// handleRegisterUser godoc
//
//	@Id				registerUser
//	@Summary		Create an account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	types.RegisterUserPayload	true	"Account details"
//	@Success		201
//	@Failure		400	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Router			/users/register [post]
func (h *Handler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	if _, err := h.store.GetUserByEmail(payload.Email); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusCreated)
}

// handleLoginUser godoc
//
//	@Id				loginUser
//	@Summary		Exchange credentials for a token
//	@Description	The returned token is sent verbatim in the Authorization header, with no "Bearer " prefix.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		types.LoginUserPayload	true	"Email and password"
//	@Success		200			{object}	types.AuthResponse
//	@Failure		400			{object}	types.ErrorResponse
//	@Failure		500			{object}	types.ErrorResponse
//	@Router			/users/login [post]
func (h *Handler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		// Deliberately identical to the wrong-password response so the reply
		// does not reveal whether an account exists.
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := auth.CreateJWT([]byte(config.Envs.JWTSecret), u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.AuthResponse{
		Token: token,
		User:  types.NewUserResponse(u),
	})
}

// handleRefreshToken godoc
//
//	@Id				refreshToken
//	@Summary		Issue a fresh token for the current session
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	types.AuthResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/refresh-token [post]
func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	u, err := h.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.CreateJWT([]byte(config.Envs.JWTSecret), u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.AuthResponse{
		Token: token,
		User:  types.NewUserResponse(u),
	})
}

// handleGetCurrentUser godoc
//
//	@Id				getCurrentUser
//	@Summary		Get the authenticated user
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	types.UserResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/me [get]
func (h *Handler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	u, err := h.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.NewUserResponse(u))
}

// handleDeleteCurrentUser godoc
//
//	@Id				deleteCurrentUser
//	@Summary		Delete the authenticated user's account
//	@Tags			users
//	@Produce		json
//	@Success		204
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/users/me [delete]
func (h *Handler) handleDeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	if err := h.store.DeleteUserById(userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteStatus(w, http.StatusNoContent)
}
