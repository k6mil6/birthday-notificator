package users

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k6mil6/birthday-notificator/internal/model"
	"log/slog"
	"time"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewStorage(db *pgxpool.Pool, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) Save(ctx context.Context, user *model.User) error {
	op := "users.Save"

	log := s.log.With(slog.String("op", op))

	query := `INSERT INTO users (id, name, birthday, email, password_hash) 
			  VALUES (@id, @name, @birthday, @email, @password_hash)`

	args := pgx.NamedArgs{
		"id":            user.ID,
		"name":          user.Name,
		"birthday":      user.Birthday,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}
		log.Error("failed to create user", "error", err)

		return err
	}

	log.Info("user created", "id", user.ID)

	return nil
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	op := "users.Get"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, name, birthday, email FROM users WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	var user dbUser
	err := s.db.QueryRow(ctx, query, args).Scan(&user.ID, &user.Name, &user.Birthday, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		log.Error("failed to get user", "error", err)

		return model.User{}, err
	}

	return model.User(user), nil
}

func (s *Storage) GetByEmail(ctx context.Context, email string) (model.User, error) {
	op := "users.GetByEmail"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, name, birthday, email, password_hash FROM users WHERE email = @email`

	args := pgx.NamedArgs{
		"email": email,
	}

	var user dbUser
	err := s.db.QueryRow(ctx, query, args).Scan(&user.ID, &user.Name, &user.Birthday, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		log.Error("failed to get user", "error", err)

		return model.User{}, err
	}

	return model.User(user), nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.User, error) {
	op := "users.GetAll"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, name, birthday, email FROM users`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Error("failed to get all users", "error", err)
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user dbUser
		err = rows.Scan(&user.ID, &user.Name, &user.Birthday, &user.Email)
		if err != nil {
			log.Error("failed to get user", "error", err)
			return nil, err
		}
		users = append(users, model.User(user))
	}
	return users, nil
}

func (s *Storage) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) error {
	op := "users.UpdateUserEmail"

	log := s.log.With(slog.String("op", op))

	query := `UPDATE users SET email = @email WHERE id = @id`

	args := pgx.NamedArgs{
		"id":    id,
		"email": email,
	}

	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExists
			}
		}
		log.Error("failed to update user email", "error", err)

		return err
	}

	return nil
}

type dbUser struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Birthday     time.Time `db:"birthday"`
	PasswordHash []byte    `db:"password_hash"`
}
