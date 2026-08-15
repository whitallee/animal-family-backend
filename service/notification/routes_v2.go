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

// RegisterV2Routes mounts the v2 notification routes.
//
// Every response has a named type rather than an ad hoc map literal, so the
// generated client knows their shape. Subscriptions are returned as
// types.PushSubscriptionResponse, which omits the p256dh and auth encryption
// keys that the stored record carries.
func (h *Handler) RegisterV2Routes(router *mux.Router) {
	router.HandleFunc("/notifications/vapid-public-key", h.handleGetVAPIDPublicKeyV2).Methods(http.MethodGet)

	router.HandleFunc("/notifications/subscribe", auth.WithJWTAuth(h.handleSubscribeV2, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/notifications/unsubscribe", auth.WithJWTAuth(h.handleUnsubscribeV2, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/notifications/subscriptions", auth.WithJWTAuth(h.handleListSubscriptions, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/notifications/test", auth.WithJWTAuth(h.handleSendTestNotification, h.userStore)).Methods(http.MethodPost)
}

// handleGetVAPIDPublicKeyV2 godoc
//
//	@Id				getVapidPublicKey
//	@Summary		Get the VAPID public key
//	@Description	Browsers need this key to create a push subscription. It is public by design.
//	@Tags			notifications
//	@Produce		json
//	@Success		200	{object}	types.VAPIDKeyResponse
//	@Router			/notifications/vapid-public-key [get]
func (h *Handler) handleGetVAPIDPublicKeyV2(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, types.VAPIDKeyResponse{
		PublicKey: config.Envs.VAPIDPublicKey,
	})
}

// handleSubscribeV2 godoc
//
//	@Id				subscribeToPush
//	@Summary		Register a push subscription
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		types.SubscribePayload	true	"Browser push subscription"
//	@Success		201				{object}	types.MessageResponse
//	@Failure		400				{object}	types.ErrorResponse
//	@Failure		403				{object}	types.ErrorResponse
//	@Failure		500				{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notifications/subscribe [post]
func (h *Handler) handleSubscribeV2(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.SubscribePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	err := h.store.CreateSubscription(types.PushSubscription{
		UserID:    userID,
		Endpoint:  payload.Endpoint,
		P256dh:    payload.Keys.P256dh,
		Auth:      payload.Keys.Auth,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, types.MessageResponse{
		Message: "subscription created successfully",
	})
}

// handleUnsubscribeV2 godoc
//
//	@Id				unsubscribeFromPush
//	@Summary		Remove a push subscription
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		types.UnsubscribePayload	true	"Endpoint to remove"
//	@Success		200				{object}	types.MessageResponse
//	@Failure		400				{object}	types.ErrorResponse
//	@Failure		403				{object}	types.ErrorResponse
//	@Failure		500				{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notifications/unsubscribe [post]
func (h *Handler) handleUnsubscribeV2(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	var payload types.UnsubscribePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validationErrors))
		return
	}

	// Scoped to the caller, so one user cannot remove another's subscription
	// by supplying their endpoint.
	if err := h.store.DeleteSubscriptionByEndpoint(userID, payload.Endpoint); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.MessageResponse{
		Message: "subscription deleted successfully",
	})
}

// handleListSubscriptions godoc
//
//	@Id				listPushSubscriptions
//	@Summary		List the caller's push subscriptions
//	@Tags			notifications
//	@Produce		json
//	@Success		200	{array}		types.PushSubscriptionResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notifications/subscriptions [get]
func (h *Handler) handleListSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	subscriptions, err := h.store.GetSubscriptionsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.NewPushSubscriptionResponses(subscriptions))
}

// handleSendTestNotification godoc
//
//	@Id				sendTestNotification
//	@Summary		Send a test push notification
//	@Description	Delivers a test message to every subscription the caller has registered, reporting the outcome for each.
//	@Tags			notifications
//	@Produce		json
//	@Success		200	{object}	types.TestNotificationResponse
//	@Failure		403	{object}	types.ErrorResponse
//	@Failure		500	{object}	types.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notifications/test [post]
func (h *Handler) handleSendTestNotification(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetuserIdFromContext(r.Context())

	subscriptions, err := h.store.GetSubscriptionsByUserId(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if len(subscriptions) == 0 {
		utils.WriteJSON(w, http.StatusOK, types.TestNotificationResponse{
			Message:           "no subscriptions found",
			SubscriptionCount: 0,
			Results:           []types.TestNotificationResult{},
		})
		return
	}

	testNotification := &types.TaskResetNotification{
		TaskName:    "Test Notification",
		TaskDesc:    "This is a test notification from your Animal Family app",
		UserID:      userID,
		SubjectName: "Test",
		SubjectType: "test",
	}

	results := make([]types.TestNotificationResult, 0, len(subscriptions))
	for _, sub := range subscriptions {
		statusCode, sendErr := h.sender.SendSingleNotificationWithStatus(sub, testNotification)

		result := types.TestNotificationResult{
			SubscriptionId: sub.SubscriptionId,
			Endpoint:       truncateEndpoint(sub.Endpoint),
			HttpStatus:     statusCode,
			Success:        sendErr == nil,
		}
		if sendErr != nil {
			result.Error = sendErr.Error()
		}

		results = append(results, result)
	}

	utils.WriteJSON(w, http.StatusOK, types.TestNotificationResponse{
		Message:           "test notification sent",
		SubscriptionCount: len(subscriptions),
		Results:           results,
	})
}

// truncateEndpoint shortens a push endpoint for display. v1 sliced the string
// directly, which panics on an endpoint shorter than the cut length.
func truncateEndpoint(endpoint string) string {
	const max = 50

	if len(endpoint) <= max {
		return endpoint
	}

	return endpoint[:max] + "..."
}
