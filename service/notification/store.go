package notification

import (
	"database/sql"
	"fmt"

	"github.com/whitallee/animal-family-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSubscription(sub types.PushSubscription) error {
	_, err := s.db.Exec(`
		INSERT INTO "pushSubscriptions" ("userId", "endpoint", "p256dh", "auth", "userAgent", "createdAt", "lastUsed")
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT ("userId", "endpoint")
		DO UPDATE SET "lastUsed" = NOW(), "userAgent" = EXCLUDED."userAgent"
	`, sub.UserID, sub.Endpoint, sub.P256dh, sub.Auth, sub.UserAgent)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSubscriptionsByUserId(userId int) ([]*types.PushSubscription, error) {
	rows, err := s.db.Query(`
		SELECT "subscriptionId", "userId", "endpoint", "p256dh", "auth", "userAgent", "createdAt", "lastUsed"
		FROM "pushSubscriptions"
		WHERE "userId" = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]*types.PushSubscription, 0)
	for rows.Next() {
		sub := new(types.PushSubscription)
		err := rows.Scan(
			&sub.SubscriptionId,
			&sub.UserID,
			&sub.Endpoint,
			&sub.P256dh,
			&sub.Auth,
			&sub.UserAgent,
			&sub.CreatedAt,
			&sub.LastUsed,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (s *Store) DeleteSubscription(subscriptionId int) error {
	result, err := s.db.Exec(`DELETE FROM "pushSubscriptions" WHERE "subscriptionId" = $1`, subscriptionId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (s *Store) DeleteSubscriptionByEndpoint(userId int, endpoint string) error {
	result, err := s.db.Exec(`DELETE FROM "pushSubscriptions" WHERE "userId" = $1 AND "endpoint" = $2`, userId, endpoint)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (s *Store) UpdateLastUsed(subscriptionId int) error {
	_, err := s.db.Exec(`UPDATE "pushSubscriptions" SET "lastUsed" = NOW() WHERE "subscriptionId" = $1`, subscriptionId)
	if err != nil {
		return err
	}

	return nil
}
