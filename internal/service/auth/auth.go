package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/k6mil6/birthday-notificator/internal/lib/jwt"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres/users"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Service struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
	secret       string
}

type UserSaver interface {
	Save(ctx context.Context, user *model.User) error
}

type UserProvider interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
	secret string,
) *Service {
	return &Service{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
		secret:       secret,
	}
}

func (a *Service) Login(ctx context.Context, login, password string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.GetByEmail(ctx, login)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Error("user not found", "login", login)

			return "", ErrInvalidCredentials
		}

		log.Error("failed to get user by login", "error", err.Error())
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", "err", err.Error())

		return "", ErrInvalidCredentials
	}

	log.Info("user logged in")

	token, err := jwt.NewToken(&user, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to create token", "error", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Service) Register(ctx context.Context, email, password, name string, birthday time.Time) (uuid.UUID, error) {
	const op = "Auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to register user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", "error", err)
		return uuid.Nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error("failed to generate id", "error", err)
		return uuid.Nil, err
	}

	user := model.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		Birthday:     birthday,
	}

	err = a.userSaver.Save(ctx, &user)
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			log.Error("user already exists", "email", email)
			return uuid.Nil, ErrUserAlreadyExists
		}
		log.Error("failed to save user", "error", err)
		return uuid.Nil, err
	}

	log.Info("user registered")

	return id, nil
}
