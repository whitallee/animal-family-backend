package notification

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
	store  types.PushSubscriptionStore
	sender *NotificationSender
}

func NewHandler(store types.PushSubscriptionStore, sender *NotificationSender) *Handler {
	return &Handler{store: store, sender: sender}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// public routes
	router.HandleFunc("/notification/vapid-public-key", h.handleGetVAPIDPublicKey).Methods(http.MethodGet)

	// authenticated routes
	router.HandleFunc("/notification/subscribe", auth.WithJWTAuth(h.handleSubscribe, nil)).Methods(http.MethodPost)
	router.HandleFunc("/notification/unsubscribe", auth.WithJWTAuth(h.handleUnsubscribe, nil)).Methods(http.MethodPost)
	router.HandleFunc("/notification/subscriptions", auth.WithJWTAuth(h.handleGetSubscriptions, nil)).Methods(http.MethodGet)
	router.HandleFunc("/notification/test", auth.WithJWTAuth(h.handleTestNotification, nil)).Methods(http.MethodPost)
}

func (h *Handler) handleGetVAPIDPublicKey(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"publicKey": config.Envs.VAPIDPublicKey,
	})
}

func (h *Handler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	// get userId from context
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var payload types.SubscribePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get user agent from request headers
	userAgent := r.Header.Get("User-Agent")

	// create subscription
	subscription := types.PushSubscription{
		UserID:    userID,
		Endpoint:  payload.Endpoint,
		P256dh:    payload.Keys.P256dh,
		Auth:      payload.Keys.Auth,
		UserAgent: userAgent,
	}

	err := h.store.CreateSubscription(subscription)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "subscription created successfully",
	})
}

func (h *Handler) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	// get userId from context
	userID := auth.GetuserIdFromContext(r.Context())

	// get JSON payload
	var payload types.UnsubscribePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// delete subscription by endpoint and userId
	err := h.store.DeleteSubscriptionByEndpoint(userID, payload.Endpoint)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "subscription deleted successfully",
	})
}

func (h *Handler) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	// get userId from context
	userID := auth.GetuserIdFromContext(r.Context())

	// get subscriptions for user
	subscriptions, err := h.store.GetSubscriptionsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, subscriptions)
}

func (h *Handler) handleTestNotification(w http.ResponseWriter, r *http.Request) {
	// get userId from context
	userID := auth.GetuserIdFromContext(r.Context())

	// create a test notification
	testNotification := &types.TaskResetNotification{
		TaskId:      0,
		TaskName:    "Test Notification",
		TaskDesc:    "This is a test notification from your Animal Family app",
		UserID:      userID,
		SubjectName: "Test",
		SubjectType: "test",
	}

	// send notification asynchronously
	go h.sender.SendTaskResetNotifications([]*types.TaskResetNotification{testNotification})

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "test notification sent",
	})
}
