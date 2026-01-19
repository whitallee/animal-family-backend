package notification

import (
	"encoding/json"
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/whitallee/animal-family-backend/types"
)

type NotificationSender struct {
	store           types.PushSubscriptionStore
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

func NewNotificationSender(store types.PushSubscriptionStore, vapidPublicKey, vapidPrivateKey, vapidSubject string) *NotificationSender {
	return &NotificationSender{
		store:           store,
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
		vapidSubject:    vapidSubject,
	}
}

func (ns *NotificationSender) SendTaskResetNotifications(tasks []*types.TaskResetNotification) {
	// Group tasks by userId for efficient subscription lookup
	tasksByUser := groupTasksByUser(tasks)

	for userId, userTasks := range tasksByUser {
		// Get all subscriptions for this user
		subscriptions, err := ns.store.GetSubscriptionsByUserId(userId)
		if err != nil {
			log.Printf("failed to get subscriptions for user %d: %v", userId, err)
			continue
		}

		// Send notification for each task to each subscription
		for _, task := range userTasks {
			for _, sub := range subscriptions {
				if err := ns.sendNotification(sub, task); err != nil {
					log.Printf("failed to send notification: %v", err)
				}
			}
		}
	}
}

func (ns *NotificationSender) sendNotification(sub *types.PushSubscription, task *types.TaskResetNotification) error {
	// Build notification payload
	payload := map[string]interface{}{
		"title": fmt.Sprintf("%s (%s)", task.TaskName, task.SubjectName),
		"body":  task.TaskDesc,
		"data": map[string]interface{}{
			"taskId": task.TaskId,
			"url":    fmt.Sprintf("/tasks/%d", task.TaskId),
		},
		"tag":                fmt.Sprintf("task-%d", task.TaskId),
		"requireInteraction": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %v", err)
	}

	// Create subscription for webpush
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	// Send notification
	resp, err := webpush.SendNotification(payloadBytes, subscription, &webpush.Options{
		VAPIDPublicKey:  ns.vapidPublicKey,
		VAPIDPrivateKey: ns.vapidPrivateKey,
		Subscriber:      ns.vapidSubject,
		TTL:             60 * 60 * 24, // 24 hours
	})

	if resp != nil {
		defer resp.Body.Close()

		// Handle 410 Gone (expired subscription) or 404 Not Found
		if resp.StatusCode == 410 || resp.StatusCode == 404 {
			log.Printf("subscription expired, deleting: %d", sub.SubscriptionId)
			ns.store.DeleteSubscription(sub.SubscriptionId)
		}
	}

	return err
}

func groupTasksByUser(tasks []*types.TaskResetNotification) map[int][]*types.TaskResetNotification {
	tasksByUser := make(map[int][]*types.TaskResetNotification)
	for _, task := range tasks {
		tasksByUser[task.UserID] = append(tasksByUser[task.UserID], task)
	}
	return tasksByUser
}
