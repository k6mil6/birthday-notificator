package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k6mil6/birthday-notificator/internal/config"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres/subscriptions"
	"github.com/k6mil6/birthday-notificator/internal/storage/postgres/users"
	"log/slog"
	"time"
)

type Storages struct {
	pgxPool       *pgxpool.Pool
	Users         *users.Storage
	Subscriptions *subscriptions.Storage
}

func NewStorages(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) (*Storages, error) {
	var pgxPool *pgxpool.Pool

	var err error
	for i := 0; i < cfg.DB.RetriesNumber; i++ {
		pgxPool, err = pgxpool.New(ctx, cfg.DB.PostgresDSN)
		if err == nil {
			break
		}
		time.Sleep(cfg.DB.RetryCooldown)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres after %d retries: %w", cfg.DB.RetriesNumber, err)
	}

	log.Info("connected to postgres", "dsn", cfg.DB.PostgresDSN)
	log.Info("applying migrations")
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DB.PostgresDSN,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	for i := 0; i < cfg.DB.RetriesNumber; i++ {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info("migrations already applied")

				break
			}

			if i < cfg.DB.RetriesNumber-1 {
				log.Info("migration attempt failed", "attempt", i+1, "error", err)
				log.Info("retrying in", "cooldown", cfg.DB.RetryCooldown.String())

				time.Sleep(cfg.DB.RetryCooldown)

				continue
			}

			return nil, fmt.Errorf("failed to apply migrations after %d retries: %w", cfg.DB.RetriesNumber, err)
		} else {
			log.Info("migrations applied")

			break
		}
	}

	return &Storages{
		pgxPool:       pgxPool,
		Users:         users.NewStorage(pgxPool, log),
		Subscriptions: subscriptions.NewStorage(pgxPool, log),
	}, nil
}

func (s *Storages) Close() {
	s.pgxPool.Close()
}
