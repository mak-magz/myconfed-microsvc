package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"google.golang.org/grpc/metadata"
)

type Handler struct {
	userv1.UnimplementedUserServiceServer
	userClient userv1.UserServiceClient
}

func NewHandler(userClient userv1.UserServiceClient) *Handler {
	return &Handler{userClient: userClient}
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	md := metadata.Pairs("x-request-id", c.Writer.Header().Get("x-request-id"))
	grpcCtx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.userClient.GetUser(grpcCtx, &userv1.GetUserRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetUser())
}
