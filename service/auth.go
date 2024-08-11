package service

import (
	"context"
	"users_service/pkg/logger"
	"users_service/storage"

	pb "users_service/genproto/users"

)

type authService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(storage storage.IStorage, log logger.ILogger) *authService {
	return &authService{
		storage: storage,
		log:     log,
	}
}

func (a *authService) Create(ctx context.Context, request *pb.CreateUser) (*pb.User, error) {

	resp, err := a.storage.Auth().Create(ctx, request)
	if err != nil {
		a.log.Error("error while creating user info in service layer", logger.Error(err))
		return &pb.User{}, err
	}

	return resp, nil
}

func (a *authService) GetByEmail(ctx context.Context, request *pb.Email) (*pb.UserByEmail, error) {

	resp, err := a.storage.Auth().GetByEmail(ctx, request)
	if err != nil {
		a.log.Error("error while getting user info by email in service layer", logger.Error(err))
		return &pb.UserByEmail{}, err
	}

	return resp, nil
}

func (a *authService) DeleteRefreshTokenByUserId(ctx context.Context, request *pb.PrimaryKey) (*pb.Void, error) {

	resp, err := a.storage.Auth().DeleteRefreshTokenByUserId(ctx, request)
	if err != nil {
		a.log.Error("error while deleting refresh token by user id in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (a *authService) StoreRefreshToken(ctx context.Context, request *pb.RefreshToken) (*pb.Void, error) {

	resp, err := a.storage.Auth().StoreRefreshToken(ctx, request)
	if err != nil {
		a.log.Error("error while storing refresh token in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (a *authService) CheckRefreshTokenExists(ctx context.Context, request *pb.RequestRefreshToken) (*pb.Void, error) {

	resp, err := a.storage.Auth().CheckRefreshTokenExists(ctx, request)
	if err != nil {
		a.log.Error("error while cheking refresh token is existing in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (a *authService) CheckEmailExists(ctx context.Context, request *pb.Email) (*pb.Void, error) {

	resp, err := a.storage.Auth().CheckEmailExists(ctx, request)
	if err != nil {
		a.log.Error("error while cheking refresh token is existing in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (a *authService) ResetPassword(ctx context.Context, request *pb.ResetPassword) (*pb.Void, error) {

	resp, err := a.storage.Auth().ResetPassword(ctx, request)
	if err != nil {
		a.log.Error("error while reseting password in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}
