package unsubscribe

import (
	"context"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"log/slog"
	"net/http"
)

type Request struct {
	UserID uuid.UUID `json:"user_id"`
}

type Response struct {
	Success bool `json:"success"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.unsubscribe.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.UserID == uuid.Nil {
			response.HandleError(w, r, http.StatusBadRequest, "user_id is required")
			return
		}

		err := userService.Unsubscribe(ctx, userID, req.UserID)
		if err != nil {
			response.HandleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		response.HandleSuccess(w, r, Response{
			Success: true,
		})
	}
}
