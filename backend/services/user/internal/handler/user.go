package handler

import (
	"context"
	"errors"
	"log/slog"

	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/domain"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailRequired),
		errors.Is(err, domain.ErrPasswordRequired),
		errors.Is(err, domain.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, repository.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, repository.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())

	default:
		// unknown / unexpected — don't leak internals to the client
		slog.Error("unhandled error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}

type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
}

type Handler struct {
	userv1.UnimplementedUserServiceServer
	svc UserService
}

func NewHandler(svc UserService) *Handler {
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

func (h *Handler) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.GetUserByEmailResponse, error) {
	slog.DebugContext(ctx, "handler GetUserByEmail", "email", req.GetEmail())

	user, err := h.svc.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.GetUserByEmailResponse{
		User: &userv1.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (h *Handler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	slog.DebugContext(ctx, "handler Register", "email", req.GetEmail())

	user, err := h.svc.CreateUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.RegisterResponse{
		User: &userv1.User{
			Id:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	slog.DebugContext(ctx, "handler: Login", "email", req.GetEmail())

	user, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.LoginResponse{
		User: &userv1.User{
			Id:    user.ID,
			Email: user.Email,
		},
		Tokens: &userv1.Tokens{
			AccessToken:  "",
			RefreshToken: "",
			ExpiresIn:    0,
		},
	}, nil
}
