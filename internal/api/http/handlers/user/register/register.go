package register

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/email"
	authservice "github.com/k6mil6/birthday-notificator/internal/service/auth"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

type Response struct {
	Success bool `json:"success"`
}

func New(ctx context.Context, log *slog.Logger, auth httpserver.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Email == "" || req.Password == "" || req.Name == "" || req.Birthday == "" {
			response.HandleError(w, r, http.StatusBadRequest, "email, password, name, birthday are required")
			return
		}

		if !email.IsValidEmail(req.Email) {
			response.HandleError(w, r, http.StatusBadRequest, "invalid email")

			return
		}

		log.Info("registering user")

		birthday, err := time.Parse("02.01.2006", req.Birthday)
		if err != nil {
			response.HandleError(w, r, http.StatusBadRequest, "incorrect birthday format should be 'DD.MM.YYYY'")

			return
		}

		id, err := auth.Register(ctx, req.Email, req.Password, req.Name, birthday)
		if err != nil {
			if errors.Is(err, authservice.ErrUserAlreadyExists) {
				response.HandleError(w, r, http.StatusConflict, "user already exists")

				return
			}
			response.HandleError(w, r, http.StatusInternalServerError, err.Error())

			return
		}

		log.Info("user registered with id", "id", id)

		response.HandleSuccess(w, r, Response{Success: true})
	}
}
