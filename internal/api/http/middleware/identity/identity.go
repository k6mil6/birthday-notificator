package identity

import (
	"context"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"log/slog"
	"net/http"
	"strings"
)

func New(log *slog.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				log.Error("failed to get id from token", "error", "Authorization header is empty")
				response.HandleError(w, r, http.StatusUnauthorized, "Authorization header is empty")

				return
			}
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				log.Error("failed to get id from token", "error", "Authorization header is not valid")
				response.HandleError(w, r, http.StatusUnauthorized, "Authorization header is not valid")

				return
			}

			if len(headerParts[1]) == 0 {
				log.Error("failed to get id from token", "error", "Authorization header is empty")
				response.HandleError(w, r, http.StatusUnauthorized, "Authorization header is empty")

				return
			}

			id, err := jwt.GetUserID(headerParts[1], secret)
			if err != nil {
				log.Error("failed to get id from token", "error", err.Error(), "token", headerParts[1])
				response.HandleError(w, r, http.StatusUnauthorized, "Failed to get id from token")

				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, jwt.KeyUserID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
