package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetUser(ctx gin.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	Register(ctx gin.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error)
}

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

func (h *Handler) Register(c *gin.Context) {
	var req userv1.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	md := metadata.Pairs("x-request-id", c.Writer.Header().Get("x-request-id"))
	grpcCtx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.userClient.Register(grpcCtx, &req)
	if err != nil {
		slog.DebugContext(c, "faled", "error", err)
		respondGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp.GetUser())
}

func httpStatusFromGRPC(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest // 400
	case codes.Unauthenticated:
		return http.StatusUnauthorized // 401
	case codes.PermissionDenied:
		return http.StatusForbidden // 403
	case codes.NotFound:
		return http.StatusNotFound // 404
	case codes.AlreadyExists:
		return http.StatusConflict // 409
	default:
		return http.StatusInternalServerError // 500
	}
}

func respondGRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(httpStatusFromGRPC(st.Code()), gin.H{"error": st.Message()})
}
