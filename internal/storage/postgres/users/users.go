package users

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *Storage) Create(ctx context.Context, user model.User) error {
	op := "users.Create"

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
		log.Error("failed to create user", "error", err)

		return err
	}

	log.Info("user created", "id", user.ID)

	return nil
}

func (s *Storage) Get(ctx context.Context, id uuid.UUID) (model.User, error) {
	op := "users.Get"

	log := s.log.With(slog.String("op", op))

	query := `SELECT id, name, birthday, email FROM users WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	var user dbUser
	err := s.db.QueryRow(ctx, query, args).Scan(&user.ID, &user.Name, &user.Birthday, &user.Email)
	if err != nil {
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

type dbUser struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Birthday     time.Time `db:"birthday"`
	PasswordHash []byte    `db:"password_hash"`
}
