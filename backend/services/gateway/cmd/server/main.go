package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/services/gateway/internal/config"
	"github.com/mak-magz/myconfed-microsvc/backend/services/gateway/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userSvcAddr = "localhost:50051"
	listenAddr  = ":8080" // Match gateway port
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("Starting gateway server...")

	config := config.Load()

	conn, err := grpc.NewClient(config.UserSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to user service", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	userClient := handler.NewHandler(userv1.NewUserServiceClient(conn))

	r := gin.Default()
	users := r.Group("/users")
	{
		users.GET("/:id", userClient.GetUser)
	}

	log.Info("Gateway server started", "addr", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Error("failed to run gateway server", "error", err)
		os.Exit(1)
	}
}
