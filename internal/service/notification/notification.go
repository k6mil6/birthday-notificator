package notification

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/lib/email/text"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"log/slog"
	"time"
)

type EmailSender interface {
	Send(email, subject, body string) error
}

type SubscriptionsProvider interface {
	GetAll(ctx context.Context) ([]model.Subscription, error)
	UpdateNotificationDate(ctx context.Context, id uuid.UUID, notificationDate time.Time) error
}

type UserProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type Service struct {
	log *slog.Logger

	scanInterval          time.Duration
	subscriptionsProvider SubscriptionsProvider
	userProvider          UserProvider
	emailSender           EmailSender
}

func New(
	log *slog.Logger,
	scanInterval time.Duration,
	subscriptionsProvider SubscriptionsProvider,
	userProvider UserProvider,
	emailSender EmailSender,
) *Service {
	return &Service{
		log:                   log,
		scanInterval:          scanInterval,
		subscriptionsProvider: subscriptionsProvider,
		userProvider:          userProvider,
		emailSender:           emailSender,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.log.Info("starting notification service")
	ticker := time.NewTicker(s.scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.log.Info("scanning subscriptions")
			s.scan(ctx)
		case <-ctx.Done():
			s.log.Info("stopping notification service")

			return
		}
	}
}

func (s *Service) scan(ctx context.Context) {
	subscriptions, err := s.subscriptionsProvider.GetAll(ctx)
	if err != nil {
		s.log.Error("failed to get subscriptions", "error", err)
		return
	}

	s.log.Info("scanned subscriptions", "count", len(subscriptions))

	for _, subscription := range subscriptions {
		currentTime := time.Now()
		if currentTime.Before(subscription.NotificationDate) {
			s.log.Info("skipping notification", "date", subscription.NotificationDate, "now", currentTime)
			continue
		}

		s.log.Info("sending notification", "date", subscription.NotificationDate)

		user, err := s.userProvider.GetByID(ctx, subscription.UserID)
		if err != nil {
			s.log.Error("failed to get user", "error", err)
			continue
		}

		subscribedAtUser, err := s.userProvider.GetByID(ctx, subscription.SubscribedAtUserID)
		if err != nil {
			s.log.Error("failed to get user", "error", err)
			continue
		}

		s.log.Info("sending notification", "user", user, "subscribed_at_user", subscribedAtUser)

		err = s.emailSender.Send(
			user.Email,
			text.BuildEmailHeader(subscribedAtUser.Name),
			text.BuildEmailBody(
				subscribedAtUser.Name,
				subscribedAtUser.Birthday,
				subscription.NotificationDate,
			),
		)
		if err != nil {
			s.log.Error("failed to send email", "error", err)
			continue
		}

		nextNotificationDate := subscription.NotificationDate.AddDate(1, 0, 0)

		err = s.subscriptionsProvider.UpdateNotificationDate(ctx, subscription.ID, nextNotificationDate)
		if err != nil {
			s.log.Error("failed to update notification date", "error", err)
			continue
		}

		s.log.Info("notification sent", "user", user, "subscribed_at_user", subscribedAtUser)
	}
}
