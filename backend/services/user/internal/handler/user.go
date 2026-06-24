package handler

import (
	"context"
	"log/slog"

	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/service"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetUser(c context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {

	slog.DebugContext(c, "user-svc GetUser", "id", req.GetId())

	h.svc.GetUser(req.GetId())

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:   req.GetId(),
			Name: "stub",
		},
	}, nil
}
