package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/pkg/logger"
	"github.com/mak-magz/myconfed-microsvc/backend/pkg/middleware"
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
	config := config.Load()

	logger := logger.Init("gateway")
	logger.Info("Starting gateway server...")

	conn, err := grpc.NewClient(config.UserSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to user service", "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Info("Closing user service connection...")
		conn.Close()
		logger.Info("User service connection closed")
	}()

	userClient := handler.NewHandler(userv1.NewUserServiceClient(conn))

	r := gin.Default()
	r.Use(middleware.RequestLogger())

	users := r.Group("/users")
	{
		users.GET("/:id", userClient.GetUser)
		users.POST("/register", userClient.Register)
	}

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting gateway service...", "port", listenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Failed to run gateway server", "error", err)
		os.Exit(1)
	case sig := <-quit:
		logger.Info("Shutting down gateway server...", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Gateway server forced to shutdown", "error", err)
		} else {
			logger.Info("Gateway server stopped")
		}
	}
}
