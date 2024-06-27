package app

import (
	"context"
	apiapp "github.com/k6mil6/birthday-notificator/internal/app/api"
	authservice "github.com/k6mil6/birthday-notificator/internal/service/auth"
	userservice "github.com/k6mil6/birthday-notificator/internal/service/user/interaction"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	API *apiapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	storages *postgres.Storages,
	tokenTTL time.Duration,
	secret string,
	port int,
) *App {
	auth := authservice.New(log, storages.Users, storages.Users, tokenTTL, secret)
	userService := userservice.New(log, storages.Subscriptions, storages.Subscriptions, storages.Users, storages.Users)

	api := apiapp.New(ctx, log, port, auth, userService, secret)

	return &App{
		API: api,
	}
}
