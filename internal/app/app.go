package app

import (
	"context"
	apiapp "github.com/k6mil6/birthday-notificator/internal/app/api"
	"github.com/k6mil6/birthday-notificator/internal/config"
	"github.com/k6mil6/birthday-notificator/internal/lib/email"
	authservice "github.com/k6mil6/birthday-notificator/internal/service/auth"
	notificationservice "github.com/k6mil6/birthday-notificator/internal/service/notification"
	userservice "github.com/k6mil6/birthday-notificator/internal/service/user/interaction"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres"
	"log/slog"
)

type App struct {
	API *apiapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	storages *postgres.Storages,
	cfg *config.Config,
) *App {
	auth := authservice.New(log, storages.Users, storages.Users, cfg.JWT.TokenTTL, cfg.JWT.Secret)
	userService := userservice.New(log, storages.Subscriptions, storages.Subscriptions, storages.Users, storages.Users)

	emailSender := email.NewSender(cfg.Email.SenderAddress, cfg.Email.SenderPassword, cfg.Email.SMTPAddress, cfg.Email.SMTPPort)

	notificationService := notificationservice.New(log, cfg.ScanInterval, storages.Subscriptions, storages.Users, emailSender)

	go notificationService.Start(ctx)

	api := apiapp.New(ctx, log, cfg.HTTPPort, auth, userService, cfg.JWT.Secret)

	return &App{
		API: api,
	}
}
