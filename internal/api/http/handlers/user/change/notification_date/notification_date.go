package notificationdate

import (
	"context"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/lib/notification/offset"
	"log/slog"
	"net/http"
)

type Request struct {
	SubscriptionID     uuid.UUID     `json:"subscription_id"`
	NotificationOffset offset.Offset `json:"notification_offset"`
}

type Response struct {
	Success bool `json:"success"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.change.notification_time.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		if req.SubscriptionID == uuid.Nil {
			response.HandleError(w, r, http.StatusBadRequest, "subscription_id is required")
			return
		}

		if req.NotificationOffset.Unit == "" {
			response.HandleError(w, r, http.StatusBadRequest, "notification_offset unit is required")
			return
		}

		if req.NotificationOffset.Value == 0 {
			response.HandleError(w, r, http.StatusBadRequest, "notification_offset value is required")
			return
		}

		notificationOffset := offset.ConvertToTimeDuration(req.NotificationOffset)

		err := userService.ChangeNotificationDate(ctx, userID, req.SubscriptionID, notificationOffset)
		if err != nil {
			response.HandleError(w, r, http.StatusInternalServerError, "could not update notification time")
			return
		}

		response.HandleSuccess(w, r, Response{
			Success: true,
		})
	}
}
