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
	store     types.PushSubscriptionStore
	userStore types.UserStore
	sender    *NotificationSender
}

func NewHandler(store types.PushSubscriptionStore, userStore types.UserStore, sender *NotificationSender) *Handler {
	return &Handler{store: store, userStore: userStore, sender: sender}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// public routes
	router.HandleFunc("/notification/vapid-public-key", h.handleGetVAPIDPublicKey).Methods(http.MethodGet)

	// authenticated routes
	router.HandleFunc("/notification/subscribe", auth.WithJWTAuth(h.handleSubscribe, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/notification/unsubscribe", auth.WithJWTAuth(h.handleUnsubscribe, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/notification/subscriptions", auth.WithJWTAuth(h.handleGetSubscriptions, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/notification/test", auth.WithJWTAuth(h.handleTestNotification, h.userStore)).Methods(http.MethodPost)
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

	// get subscriptions for user
	subscriptions, err := h.store.GetSubscriptionsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get subscriptions: %v", err))
		return
	}

	if len(subscriptions) == 0 {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message":        "no subscriptions found",
			"subscriptionCount": 0,
		})
		return
	}

	// create a test notification
	testNotification := &types.TaskResetNotification{
		TaskId:      0,
		TaskName:    "Test Notification",
		TaskDesc:    "This is a test notification from your Animal Family app",
		UserID:      userID,
		SubjectName: "Test",
		SubjectType: "test",
	}

	// send notification synchronously for testing and capture any errors
	results := []map[string]interface{}{}
	for _, sub := range subscriptions {
		result := map[string]interface{}{
			"subscriptionId": sub.SubscriptionId,
			"endpoint":       sub.Endpoint[:50] + "...", // truncate for readability
		}
		if err := h.sender.SendSingleNotification(sub, testNotification); err != nil {
			result["success"] = false
			result["error"] = err.Error()
		} else {
			result["success"] = true
		}
		results = append(results, result)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":          "test notification sent",
		"subscriptionCount": len(subscriptions),
		"results":          results,
	})
}
