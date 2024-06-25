package subscriptions

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"log/slog"
	"time"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewStorage(db *pgxpool.Pool, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) Create(ctx context.Context, subscription model.Subscription) error {
	op := "subscriptions.Create"

	log := s.log.With(slog.String("op", op))

	query := `INSERT INTO subscriptions (id, user_id, subscribed_at_user_id, notification_time) 
			  VALUES (@id, @user_id, @subscribed_at_user_id, @notification_time)`

	args := pgx.NamedArgs{
		"id":                    subscription.ID,
		"user_id":               subscription.UserID,
		"subscribed_at_user_id": subscription.SubscribedAtUserID,
		"notification_time":     subscription.NotificationTime,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("failed to create subscription", "error", err)
		return err
	}

	log.Info("subscription created", "id", subscription.ID)

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	op := "subscriptions.Delete"

	log := s.log.With(slog.String("op", op))

	query := `DELETE FROM subscriptions WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("failed to delete subscription", "error", err)
		return err
	}

	log.Info("subscription deleted", "id", id)

	return nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Subscription, error) {
	op := "subscriptions.GetAll"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, user_id, subscribed_at_user_id, notification_time FROM subscriptions`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Error("failed to get subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]model.Subscription, 0)
	for rows.Next() {
		var subscription dbSubscription

		err = rows.Scan(&subscription.ID, &subscription.UserID, &subscription.SubscribedAtUserID, &subscription.NotificationTime)
		if err != nil {
			log.Error("failed to get subscription", "error", err)

			return nil, err
		}

		subscriptions = append(subscriptions, model.Subscription(subscription))
	}

	return subscriptions, nil
}

type dbSubscription struct {
	ID                 uuid.UUID `db:"id"`
	UserID             uuid.UUID `db:"user_id"`
	SubscribedAtUserID uuid.UUID `db:"subscribed_at_user_id"`
	NotificationTime   time.Time `db:"notification_time"`
}
