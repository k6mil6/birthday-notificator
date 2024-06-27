package subscribe

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/lib/notification/offset"
	"github.com/k6mil6/birthday-notificator/internal/service/user/interaction"
	"log/slog"
	"net/http"
)

type Request struct {
	UserID             uuid.UUID     `json:"user_id"`
	NotificationOffset offset.Offset `json:"notification_offset"`
}

type Response struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.subscribe.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		if req.UserID == uuid.Nil {
			response.HandleError(w, r, http.StatusBadRequest, "user_id is required")
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

		id, err := userService.Subscribe(ctx, userID, req.UserID, notificationOffset)
		if err != nil {
			if errors.Is(err, interaction.ErrSubscriptionAlreadyExists) {
				response.HandleError(w, r, http.StatusConflict, "user already subscribed")
				return
			}
			response.HandleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		response.HandleSuccess(w, r, Response{
			SubscriptionID: id,
		})
	}
}
