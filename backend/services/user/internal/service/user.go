package service

import (
	"context"
	"fmt"

	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, id string) {
	fmt.Println("service: GetUser", id)
	s.repo.GetUser(ctx, id)
}
