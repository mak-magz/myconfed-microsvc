package handler

import (
	"context"
	"errors"
	"log/slog"

	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/domain"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailRequired),
		errors.Is(err, domain.ErrPasswordRequired):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, repository.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		// unknown / unexpected — don't leak internals to the client
		slog.Error("unhandled error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}

type UserService interface {
	GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	CreateUser(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error)
}

type Handler struct {
	userv1.UnimplementedUserServiceServer
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetUser(c context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	slog.DebugContext(c, "handler GetUser", "id", req.GetId())

	user, err := h.svc.GetUser(c, req.GetId())

	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}
