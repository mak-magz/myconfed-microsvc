package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	GetUser(ctx context.Context, id string) (*User, error)
}

type DBRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &DBRepository{db: db}
}

func (r *DBRepository) GetUser(ctx context.Context, id string) (*User, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "user-repo GetUser", "latency", time.Since(start))
	}()

	slog.DebugContext(ctx, "user-repo GetUser", "id", id)

	var user User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name, created_at, updated_at 
		FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Email:     "stub@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
