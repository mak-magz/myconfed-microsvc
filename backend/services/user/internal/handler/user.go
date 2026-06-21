package handler

import (
	"fmt"

	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetUser(id string) {
	fmt.Println("handler: GetUser", id)
	h.svc.GetUser(id)
}
