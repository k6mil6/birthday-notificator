package subscriptions

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	query := `INSERT INTO subscriptions (id, user_id, subscribed_at_user_id, notification_date) 
			  VALUES (@id, @user_id, @subscribed_at_user_id, @notification_date)`

	args := pgx.NamedArgs{
		"id":                    subscription.ID,
		"user_id":               subscription.UserID,
		"subscribed_at_user_id": subscription.SubscribedAtUserID,
		"notification_date":     subscription.NotificationDate,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrSubscriptionAlreadyExists
			}
		}
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

func (s *Storage) GetByUsersIDs(ctx context.Context, userID, subscribedAtUserID uuid.UUID) (model.Subscription, error) {
	op := "subscriptions.GetBySubscribedAtUserID"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, user_id, subscribed_at_user_id, notification_date FROM subscriptions 
              WHERE subscribed_at_user_id = @subscribed_at_user_id AND user_id = @user_id`

	args := pgx.NamedArgs{
		"subscribed_at_user_id": subscribedAtUserID,
		"user_id":               userID,
	}

	var subscription dbSubscription
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.SubscribedAtUserID,
		&subscription.NotificationDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, ErrSubscriptionNotFound
		}
		log.Error("failed to get subscription", "error", err)

		return model.Subscription{}, err
	}

	return model.Subscription(subscription), nil
}

func (s *Storage) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	op := "subscriptions.GetAll"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, user_id, subscribed_at_user_id, notification_date FROM subscriptions 
              WHERE user_id = @user_id`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		log.Error("failed to get subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]model.Subscription, 0)
	for rows.Next() {
		var subscription dbSubscription

		err = rows.Scan(&subscription.ID, &subscription.UserID, &subscription.SubscribedAtUserID, &subscription.NotificationDate)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrSubscriptionNotFound
			}
			log.Error("failed to get subscription", "error", err)

			return nil, err
		}

		subscriptions = append(subscriptions, model.Subscription(subscription))
	}

	return subscriptions, nil
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	op := "subscriptions.GetByID"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, user_id, subscribed_at_user_id, notification_date FROM subscriptions
			  WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	var subscription dbSubscription
	err := s.db.QueryRow(ctx, query, args).Scan(&subscription.ID, &subscription.UserID, &subscription.SubscribedAtUserID, &subscription.NotificationDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, ErrSubscriptionNotFound
		}
		log.Error("failed to get subscription", "error", err)

		return model.Subscription{}, err
	}

	return model.Subscription(subscription), nil
}

func (s *Storage) UpdateNotificationDate(ctx context.Context, id uuid.UUID, notificationDate time.Time) error {
	op := "subscriptions.UpdateNotificationDate"

	log := s.log.With(slog.String("op", op))

	query := `UPDATE subscriptions SET notification_date = @notification_date WHERE id = @id`

	args := pgx.NamedArgs{
		"notification_date": notificationDate,
		"id":                id,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		log.Error("failed to update subscription", "error", err)
		return err
	}

	log.Info("subscription updated", "id", id)

	return nil
}

type dbSubscription struct {
	ID                 uuid.UUID `db:"id"`
	UserID             uuid.UUID `db:"user_id"`
	SubscribedAtUserID uuid.UUID `db:"subscribed_at_user_id"`
	NotificationDate   time.Time `db:"notification_date"`
}
