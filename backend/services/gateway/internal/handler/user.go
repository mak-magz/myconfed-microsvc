package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
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

	resp, err := h.userClient.GetUser(c, &userv1.GetUserRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.GetUser())
}
