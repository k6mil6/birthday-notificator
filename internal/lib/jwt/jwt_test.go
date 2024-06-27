package jwt_test

import (
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"testing"
	"time"
)

type jwtTestCase struct {
	user    model.User
	secret  string
	expired bool
}

var jwtTestCases = []jwtTestCase{
	{
		user:    model.User{ID: uuid.New(), Email: "test1@mail.com", PasswordHash: []byte("hash")},
		secret:  "secret",
		expired: false,
	},
	{
		user:    model.User{ID: uuid.New(), Email: "test2@mail.com", PasswordHash: []byte("hash")},
		secret:  "secret",
		expired: true,
	},
	{
		user:    model.User{ID: uuid.New(), Email: "test3@mail.com", PasswordHash: []byte("hash")},
		secret:  "nosoeaod",
		expired: false,
	},
	{
		user:    model.User{ID: uuid.New(), Email: "test4@mail.com", PasswordHash: []byte("hash")},
		secret:  "nosoeaod",
		expired: true,
	},
}

func TestGetUserID(t *testing.T) {
	for _, tc := range jwtTestCases {
		t.Run(tc.user.Email, func(t *testing.T) {
			duration := time.Hour
			if tc.expired {
				duration = time.Nanosecond

				token, err := jwt.NewToken(&tc.user, duration, tc.secret)
				if err != nil {
					t.Errorf("expected nil, got %v", err)
				}

				id, err := jwt.GetUserID(token, tc.secret)
				if err == nil {
					t.Errorf("expected error, got res %v", id)
				}

				return
			}

			token, err := jwt.NewToken(&tc.user, duration, tc.secret)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			id, err := jwt.GetUserID(token, tc.secret)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			if id != tc.user.ID {
				t.Errorf("expected %v, got %v", tc.user.ID, id)
			}
		})
	}
}
