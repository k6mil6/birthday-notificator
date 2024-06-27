package login

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	authservice "github.com/k6mil6/birthday-notificator/internal/service/auth"
	"log/slog"
	"net/http"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
}

func New(ctx context.Context, log *slog.Logger, auth httpserver.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.user.login.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			response.HandleError(w, r, http.StatusBadRequest, err.Error())

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Email == "" || req.Password == "" {
			response.HandleError(w, r, http.StatusBadRequest, "email and password are required")

			return
		}

		token, err := auth.Login(ctx, req.Email, req.Password)
		if err != nil {
			if errors.Is(err, authservice.ErrInvalidCredentials) {
				response.HandleError(w, r, http.StatusUnauthorized, "invalid credentials")

				return
			}

			response.HandleError(w, r, http.StatusInternalServerError, "internal server error")

			return
		}

		response.HandleSuccess(w, r, Response{Token: token})
	}
}
