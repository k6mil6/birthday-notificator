package notification

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"log/slog"
	"time"
)

type SubscriptionsProvider interface {
	GetAll(ctx context.Context) ([]model.Subscription, error)
}

type UserProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type Service struct {
	log *slog.Logger

	scanInterval          time.Duration
	subscriptionsProvider SubscriptionsProvider
	userProvider          UserProvider
}

func New(log *slog.Logger, scanInterval time.Duration, subscriptionsProvider SubscriptionsProvider, userProvider UserProvider) *Service {
	return &Service{
		log:                   log,
		scanInterval:          scanInterval,
		subscriptionsProvider: subscriptionsProvider,
		userProvider:          userProvider,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.log.Info("starting notification service")
	ticker := time.NewTicker(s.scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

	for _, subscription := range subscriptions {
		if subscription.NotificationDate.After(time.Now()) || subscription.NotificationDate == time.Now() {
			user, err := s.userProvider.GetByID(ctx, subscription.UserID)
			if err != nil {
				s.log.Error("failed to get user", "error", err)
				continue
			}

			s.log.Info("sending notification", "user", user)

			// TODO: send notification to email and update notification time
		}
	}
}
