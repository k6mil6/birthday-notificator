package subscriptions

import (
	"context"
	"errors"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/model"
	subscriptionsservice "github.com/k6mil6/birthday-notificator/internal/storage/postgres/subscriptions"
	"log/slog"
	"net/http"
	"time"
)

type Response struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	SubscribedAtUserID uuid.UUID `json:"subscribed_at_user_id"`
	NotificationDate   time.Time `json:"notification_date"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.subscriptions.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		subscriptions, err := userService.GetAllSubscriptions(ctx, userID)
		if err != nil {
			if errors.Is(err, subscriptionsservice.ErrSubscriptionNotFound) {
				response.HandleError(w, r, http.StatusNotFound, "no subscriptions found")
				return
			}
			log.Error("failed to get subscriptions", "error", err.Error())
			response.HandleError(w, r, http.StatusInternalServerError, "failed to get subscriptions")
			return
		}

		response.HandleSuccess(w, r, Response{
			Subscriptions: subscriptionsToResponse(subscriptions),
		})
	}
}

func subscriptionsToResponse(subscriptions []model.Subscription) []Subscription {
	responseSubscriptions := make([]Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		responseSubscriptions = append(responseSubscriptions, Subscription{
			ID:                 subscription.ID,
			UserID:             subscription.UserID,
			SubscribedAtUserID: subscription.SubscribedAtUserID,
			NotificationDate:   subscription.NotificationDate,
		})
	}

	return responseSubscriptions
}
