package service

import (
	"context"
	pb "users_service/genproto/users"
	"users_service/pkg/logger"
	"users_service/storage"
)

type IServiceManager interface {
	AuthService() IAuthService
	UsersService() IUsersService
}

type IAuthService interface {
	Create(context.Context, *pb.CreateUser) (*pb.User, error)
	GetByEmail(context.Context, *pb.Email) (*pb.User, error)
	DeleteRefreshTokenByUserId(context.Context, *pb.PrimaryKey) (*pb.Void, error)
	StoreRefreshToken(context.Context, *pb.RefreshToken) (*pb.Void, error)
	CheckRefreshTokenExists(context.Context, *pb.RequestRefreshToken) (*pb.Void, error)
	CheckEmailExists(context.Context, *pb.Email) (*pb.Void, error)
	ResetPassword(context.Context, *pb.ResetPassword) (*pb.Void, error)
}

type IUsersService interface {
	GetById(context.Context, *pb.PrimaryKey) (*pb.User, error)
	GetAll(context.Context, *pb.GetListRequest) (*pb.Users, error)
	Update(context.Context, *pb.UpdateUser) (*pb.UpdatedUser, error)
	Delete(context.Context, *pb.PrimaryKey) (*pb.Void, error)
	ChangePassword(context.Context, *pb.ChangePassword) (*pb.Void, error)
	ChangeUserRole(context.Context, *pb.ChangeUserRole) (*pb.Void, error)
}

type ServiceManager struct {
	authService  IAuthService
	usersService IUsersService
}

func NewServiceManager(storage storage.IStorage, log logger.ILogger) IServiceManager {
	return &ServiceManager{
		authService:  NewAuthService(storage, log),
		usersService: NewUsersService(storage, log),
	}
}

func (s *ServiceManager) AuthService() IAuthService {
	return s.authService
}

func (s *ServiceManager) UsersService() IUsersService {
	return s.usersService
}
