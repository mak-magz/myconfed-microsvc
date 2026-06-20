package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listenAddr = ":50052" // Match gateway port
)

type server struct {
	// pb.UnimplementedGatewayServer
}

func main() {
	ctx := context.Background()

	serverImpl := &server{}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	// pb.RegisterGatewayServer(grpcServer, serverImpl)

	// Start listening
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	println("Gateway listening on", listenAddr)
	if err = grpcServer.Serve(listener); err != nil {
		panic("failed to serve: " + err.Error())
	}
}
