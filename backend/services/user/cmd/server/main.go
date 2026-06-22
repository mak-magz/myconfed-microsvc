package main

import (
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	userv1 "github.com/mak-magz/myconfed-microsvc/backend/gen/user/v1"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/config"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/handler"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/repository"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	config := config.Load()

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	defer db.Close()

	// wiring: repository -> service -> handler
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	hnd := handler.NewHandler(svc)
	log.Println("starting user service . . .")

	listen, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		log.Fatal("failed to listen", "error", err)
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, hnd)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal("failed to serve", "error", err)
	}
}
