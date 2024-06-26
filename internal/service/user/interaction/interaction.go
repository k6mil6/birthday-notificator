package interaction

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/lib/birthday"
	"github.com/k6mil6/birthday-notificator/internal/model"
	subscriptionsstorage "github.com/k6mil6/birthday-notificator/internal/storage/postgres/subscriptions"
	"log/slog"
	"time"
)

type SubscriptionSaver interface {
	Create(ctx context.Context, subscription model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateNotificationDate(ctx context.Context, id uuid.UUID, notificationDate time.Time) error
}

type SubscriptionProvider interface {
	GetByUsersIDs(ctx context.Context, userID, subscribedAtUserID uuid.UUID) (model.Subscription, error)
	GetAllForUser(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error)
}

type UserProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
}

type UserSaver interface {
	UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) error
}

type Service struct {
	log *slog.Logger

	subscriptionSaver    SubscriptionSaver
	subscriptionProvider SubscriptionProvider
	userProvider         UserProvider
	userSaver            UserSaver
}

func New(
	log *slog.Logger,
	subscriptionSaver SubscriptionSaver,
	subscriptionProvider SubscriptionProvider,
	userProvider UserProvider,
	userSaver UserSaver,
) *Service {
	return &Service{
		log:                  log,
		subscriptionSaver:    subscriptionSaver,
		subscriptionProvider: subscriptionProvider,
		userProvider:         userProvider,
		userSaver:            userSaver,
	}
}

func (s *Service) Subscribe(
	ctx context.Context,
	userID,
	subscribedAtUserID uuid.UUID,
	notificationOffset time.Duration,
) (uuid.UUID, error) {
	op := "interaction.Subscribe"

	log := s.log.With(slog.String("op", op))

	log.Info("attempting to subscribe user", "user_id", userID, "subscribed_at_user_id", subscribedAtUserID)

	subscribedAtUser, err := s.userProvider.GetByID(ctx, subscribedAtUserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return uuid.Nil, err
	}

	nextBirthday := birthday.CalculateNextBirthday(subscribedAtUser.Birthday)
	notificationDate := nextBirthday.Add(-notificationOffset)

	log.Info("user is being subscribed", "user_id", userID, "subscribed_at_user_id", subscribedAtUserID, "notification_date", notificationDate)

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error("failed to generate subscription id", "error", err)
		return uuid.Nil, err
	}

	subscription := model.Subscription{
		ID:                 id,
		UserID:             userID,
		SubscribedAtUserID: subscribedAtUserID,
		NotificationDate:   notificationDate,
	}

	err = s.subscriptionSaver.Create(ctx, subscription)
	if err != nil {
		log.Error("failed to create subscription", "error", err)
		return uuid.Nil, err
	}

	log.Info("user subscribed", "user_id", userID, "subscribed_at_user_id", subscribedAtUserID, "notification_date", notificationDate)

	return id, nil
}

func (s *Service) Unsubscribe(ctx context.Context, userID, subscribedAtUserID uuid.UUID) error {
	op := "interaction.Unsubscribe"

	log := s.log.With(slog.String("op", op))

	log.Info("attempting to unsubscribe user", "user_id", userID, "subscribed_at_user_id", subscribedAtUserID)

	subscription, err := s.subscriptionProvider.GetByUsersIDs(ctx, userID, subscribedAtUserID)
	if err != nil {
		if errors.Is(err, subscriptionsstorage.ErrSubscriptionNotFound) {
			log.Info("user is not subscribed", "user_id", userID, "subscribed_at_user_id", subscribedAtUserID)
			return nil
		}
		log.Error("failed to get subscription", "error", err)
		return err
	}

	log.Info(
		"user is being unsubscribed",
		"user_id", userID,
		"subscribed_at_user_id", subscribedAtUserID,
		"notification_date", subscription.NotificationDate,
	)

	err = s.subscriptionSaver.Delete(ctx, subscription.ID)
	if err != nil {
		log.Error("failed to delete subscription", "error", err)
		return err
	}

	log.Info(
		"user unsubscribed",
		"user_id", userID,
		"subscribed_at_user_id", subscribedAtUserID,
		"notification_date", subscription.NotificationDate,
	)

	return nil
}

func (s *Service) GetAllSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	op := "interaction.GetAllSubscriptions"

	log := s.log.With(slog.String("op", op))

	log.Info("attempting to get all subscriptions")

	subscriptions, err := s.subscriptionProvider.GetAllForUser(ctx, userID)
	if err != nil {
		log.Error("failed to get subscriptions", "error", err)
		return nil, err
	}

	log.Info("got all subscriptions")

	return subscriptions, nil
}

func (s *Service) ChangeNotificationDate(ctx context.Context, userID, subscriptionID uuid.UUID, notificationOffset time.Duration) error {
	op := "interaction.ChangeNotificationDate"

	log := s.log.With(slog.String("op", op))

	log.Info(
		"attempting to change notification date",
		"user_id", userID,
		"subscription_id", subscriptionID,
		"notification_offset", notificationOffset,
	)

	subscription, err := s.subscriptionProvider.GetByID(ctx, subscriptionID)
	if err != nil {
		log.Error("failed to get subscription", "error", err)
		return err
	}

	if subscription.UserID != userID {
		log.Error("user is not subscribed", "user_id", userID, "subscribed_at_user_id", subscription.SubscribedAtUserID)
		return ErrNotAllowed
	}

	subscribedAtUser, err := s.userProvider.GetByID(ctx, subscription.SubscribedAtUserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return err
	}

	nextBirthday := birthday.CalculateNextBirthday(subscribedAtUser.Birthday)
	notificationDate := nextBirthday.Add(-notificationOffset)

	err = s.subscriptionSaver.UpdateNotificationDate(ctx, subscriptionID, notificationDate)
	if err != nil {
		log.Error("failed to update subscription", "error", err)
		return err
	}

	log.Info(
		"notification date changed",
		"user_id", userID,
		"subscription_id", subscriptionID,
		"notification_date", subscription.NotificationDate,
	)

	return nil
}

func (s *Service) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error {
	op := "interaction.UpdateUserEmail"

	log := s.log.With(slog.String("op", op))

	log.Info(
		"attempting to update user email",
		"user_id", userID,
		"email", email,
	)

	err := s.userSaver.UpdateUserEmail(ctx, userID, email)
	if err != nil {
		log.Error("failed to update user", "error", err)
		return err
	}

	log.Info(
		"user email updated",
		"user_id", userID,
		"email", email,
	)

	return nil
}

func (s *Service) GetAllUsers(ctx context.Context) ([]model.User, error) {
	op := "interaction.GetAllUsers"

	log := s.log.With(slog.String("op", op))

	log.Info("attempting to get all users")

	users, err := s.userProvider.GetAll(ctx)
	if err != nil {
		log.Error("failed to get users", "error", err)
		return nil, err
	}

	log.Info("got all users")

	return users, nil
}
