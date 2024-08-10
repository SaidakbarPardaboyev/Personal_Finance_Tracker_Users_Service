package service

import (
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"

	pb "auth_service/genproto/users"
)

type userService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewUsersService(storage storage.IStorage, log logger.ILogger) *userService {
	return &userService{
		storage: storage,
		log:     log,
	}
}

func (u *userService) GetUserProfile(ctx context.Context, request *pb.PrimaryKey) (*pb.User, error) {

	resp, err := u.storage.Users().GetUserProfile(ctx, request)
	if err != nil {
		u.log.Error("error while getting user info in service layer", logger.Error(err))
		return &pb.User{}, err
	}

	return resp, nil
}

func (u *userService) UpdateUserProfile(ctx context.Context, request *pb.UpdateUser) (*pb.UpdateProfileResponce, error) {

	resp, err := u.storage.Users().UpdateUserProfile(ctx, request)
	if err != nil {
		u.log.Error("error while updating user info in service layer", logger.Error(err))
		return &pb.UpdateProfileResponce{}, err
	}

	return resp, nil
}

func (u *userService) ChangePassword(ctx context.Context, request *pb.ChangePassword) (*pb.Message, error) {

	iscurrent, err := u.storage.Users().CheckPasswordExisis(ctx, request)
	if err != nil {
		u.log.Error("error while checking current password is currect in service layer", logger.Error(err))
		return &pb.Message{}, err
	}

	if !iscurrent {
		u.log.Error("error while current password is not correct in service layer", logger.Error(err))
		return &pb.Message{}, err
	}

	resp, err := u.storage.Users().ChangePassword(ctx, request)
	if err != nil {
		u.log.Error("error while changing password in service layer", logger.Error(err))
		return &pb.Message{}, err
	}

	return resp, nil
}
