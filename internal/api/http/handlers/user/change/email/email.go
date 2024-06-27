package email

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/email"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/service/user/interaction"
	"log/slog"
	"net/http"
)

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	Success bool `json:"success"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.change.email.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, "error decoding request body")

			return
		}

		if req.Email == "" {
			response.HandleError(w, r, http.StatusBadRequest, "email is required")

			return
		}

		if !email.IsValidEmail(req.Email) {
			response.HandleError(w, r, http.StatusBadRequest, "invalid email")

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err := userService.UpdateUserEmail(ctx, userID, req.Email)
		if err != nil {
			if errors.Is(err, interaction.ErrEmailAlreadyExists) {
				response.HandleError(w, r, http.StatusConflict, "user with this email already exists")

				return
			}
			response.HandleError(w, r, http.StatusInternalServerError, "could not update email")

			return
		}

		response.HandleSuccess(w, r, Response{
			Success: true,
		})
	}
}
