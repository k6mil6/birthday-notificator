package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"time"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password, name string, birthday time.Time) (uuid.UUID, error)
}

type UserService interface {
	Subscribe(ctx context.Context, userID, subscribedAtUserID uuid.UUID, notificationOffset time.Duration) (uuid.UUID, error)
	Unsubscribe(ctx context.Context, userID, subscribedAtUserID uuid.UUID) error
	GetAllSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error)
	ChangeNotificationDate(ctx context.Context, userID, subscriptionID uuid.UUID, notificationOffset time.Duration) error
	UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	IsUserExists(ctx context.Context, userID uuid.UUID) (bool, error)
}
