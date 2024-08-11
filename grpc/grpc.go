package grpc

import (
	pb "users_service/genproto/users"
	"users_service/pkg/logger"
	"users_service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(services service.IServiceManager, log logger.ILogger) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, services.AuthService())
	pb.RegisterUsersServiceServer(grpcServer, services.UsersService())

	reflection.Register(grpcServer)
	return grpcServer
}
