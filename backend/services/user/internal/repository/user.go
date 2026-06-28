package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/domain"
)

type Repository interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type DBRepository struct {
	db *sqlx.DB
}

type User struct {
	ID        string     `db:"id"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewRepository(db *sqlx.DB) Repository {
	return &DBRepository{db: db}
}

func (r *DBRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "repo GetUser", "latency", time.Since(start))
	}()

	slog.DebugContext(ctx, "repo GetUser", "id", id)

	var user User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, created_at, updated_at FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (r *DBRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "repo CreateUser", "latency", time.Since(start))
	}()

	userRow := User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
	}

	slog.DebugContext(ctx, "repo CreateUser", "user", userRow)

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		userRow.ID, userRow.Email, userRow.Password, userRow.CreatedAt, nil)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *DBRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "repo GetUserByEmail", "latency", time.Since(start))
	}()

	slog.DebugContext(ctx, "repo GetUserByEmail", "email", email)

	var user User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, created_at, updated_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
