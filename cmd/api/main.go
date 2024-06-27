package main

import (
	"context"
	"github.com/k6mil6/birthday-notificator/internal/app"
	"github.com/k6mil6/birthday-notificator/internal/config"
	"github.com/k6mil6/birthday-notificator/internal/lib/logger"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Info("connecting to db", slog.String("dsn", cfg.DB.PostgresDSN))
	storages, err := postgres.NewStorages(ctx, log, cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err.Error())
		return
	}

	log.Info("connected to db", slog.String("dsn", cfg.DB.PostgresDSN))

	defer storages.Close()

	application := app.New(ctx, log, storages, cfg)

	go func() {
		application.API.MustRun()
	}()

	log.Info("started server")

	<-ctx.Done()

	log.Info("shutting down")
}
