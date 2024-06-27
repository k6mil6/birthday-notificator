package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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
	postgresConnectionString string,
	maxRetries int,
	retryCooldown time.Duration,
) (*Storages, error) {
	var pgxPool *pgxpool.Pool

	var err error
	for i := 0; i < maxRetries; i++ {
		pgxPool, err = pgxpool.New(ctx, postgresConnectionString)
		if err == nil {
			break
		}
		time.Sleep(retryCooldown)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres after %d retries: %w", maxRetries, err)
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
