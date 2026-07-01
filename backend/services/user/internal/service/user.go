package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/domain"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, id string) (*domain.User, error) {
	slog.InfoContext(ctx, "service: GetUser", "id", id)
	user, err := s.repo.GetUserById(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	slog.InfoContext(ctx, "service: GetUserByEmail", "email", email)
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	slog.InfoContext(ctx, "service: CreateUser", "email", email)

	user, err := domain.NewUser(uuid.NewString(), email, password)
	if err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "service: failed to hash password", "error", err)
		return nil, err
	}

	user.Password = string(hashed)

	user, err = s.repo.CreateUser(ctx, user)

	if err != nil {
		slog.ErrorContext(ctx, "service: failed to persist user", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.User, error) {
	slog.InfoContext(ctx, "service: Login", "email", email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			slog.ErrorContext(ctx, "service: user not found", "error", err)
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			slog.ErrorContext(ctx, "service: password mismatch", "error", err)
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}
