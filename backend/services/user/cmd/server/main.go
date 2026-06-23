package main

import (
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/pkg/logger"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/config"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/handler"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := config.Load()

	logger := logger.Init("user")
	logger.Info("Starting user service...")

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	repo := repository.NewRepository()
	svc := service.NewService(repo)
	hnd := handler.NewHandler(svc)

	listen, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, hnd)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listen); err != nil {
		logger.Error("failed to serve", "error", err)
		os.Exit(1)
	}

	logger.Info("User service started successfully", "port", config.GrpcPort)
}
