package grpc

import (
	"users_service/pkg/logger"
	"users_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(storage *storage.IStorage, log logger.ILogger) *grpc.Server {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	return grpcServer
}
