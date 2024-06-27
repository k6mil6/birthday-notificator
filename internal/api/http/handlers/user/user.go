package user

import (
	"context"
	"github.com/google/uuid"
	httpserver "github.com/k6mil6/birthday-notificator/internal/api/http"
	"github.com/k6mil6/birthday-notificator/internal/api/http/response"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"log/slog"
	"net/http"
)

type Response struct {
	Users []ResponseUser `json:"users"`
}

type ResponseUser struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Birthday string    `json:"birthday"`
}

func New(ctx context.Context, log *slog.Logger, userService httpserver.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := r.Context().Value(jwt.KeyUserID).(uuid.UUID)

		exists, err := userService.IsUserExists(ctx, userID)
		if err != nil {
			log.Error("failed to get user", "error", err.Error())
			response.HandleError(w, r, http.StatusUnauthorized, "failed to get user")
			return
		}

		if !exists {
			log.Error("user not found", "error", err.Error())
			response.HandleError(w, r, http.StatusUnauthorized, "user not found")
			return
		}

		users, err := userService.GetAllUsers(ctx)
		if err != nil {
			log.Error("failed to get users", "err", err.Error())
			response.HandleError(w, r, http.StatusInternalServerError, "failed to get users")
			return
		}

		response.HandleSuccess(w, r, Response{Users: usersToResponse(users)})
	}
}

func usersToResponse(users []model.User) []ResponseUser {
	responseUsers := make([]ResponseUser, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, ResponseUser{
			ID:       user.ID,
			Name:     user.Name,
			Birthday: user.Birthday.Format("02.01.2006"),
		})
	}
	return responseUsers
}
