package service

import (
	"context"
	"users_service/pkg/logger"
	"users_service/storage"

	pb "users_service/genproto/users"
)

type userService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedUsersServiceServer
}

func NewUsersService(storage storage.IStorage, log logger.ILogger) *userService {
	return &userService{
		storage: storage,
		log:     log,
	}
}

func (u *userService) GetById(ctx context.Context, request *pb.PrimaryKey) (*pb.User, error) {

	resp, err := u.storage.Users().GetById(ctx, request)
	if err != nil {
		u.log.Error("error while getting user info in service layer", logger.Error(err))
		return &pb.User{}, err
	}

	return resp, nil
}

func (u *userService) GetAll(ctx context.Context, request *pb.GetListRequest) (*pb.Users, error) {

	resp, err := u.storage.Users().GetAll(ctx, request)
	if err != nil {
		u.log.Error("error while getting all users info in service layer", logger.Error(err))
		return &pb.Users{}, err
	}

	return resp, nil
}

func (u *userService) Update(ctx context.Context, request *pb.UpdateUser) (*pb.UpdatedUser, error) {

	resp, err := u.storage.Users().Update(ctx, request)
	if err != nil {
		u.log.Error("error while updating user info in service layer", logger.Error(err))
		return &pb.UpdatedUser{}, err
	}

	return resp, nil
}

func (u *userService) Delete(ctx context.Context, request *pb.PrimaryKey) (*pb.Void, error) {

	resp, err := u.storage.Users().Delete(ctx, request)
	if err != nil {
		u.log.Error("error while deleting user info in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (u *userService) ChangePassword(ctx context.Context, request *pb.ChangePassword) (*pb.Void, error) {

	iscurrent, err := u.storage.Users().CheckPasswordExisis(ctx, request)
	if err != nil {
		u.log.Error("error while checking current password is currect in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	if !iscurrent {
		u.log.Error("error while current password is not correct in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	resp, err := u.storage.Users().ChangePassword(ctx, request)
	if err != nil {
		u.log.Error("error while changing password in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}

func (u *userService) ChangeUserRole(ctx context.Context, request *pb.ChangeUserRole) (*pb.Void, error) {

	resp, err := u.storage.Users().ChangeUserRole(ctx, request)
	if err != nil {
		u.log.Error("error while changing user role in service layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return resp, nil
}
