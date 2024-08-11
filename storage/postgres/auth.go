package postgres

import (
	"context"
	"fmt"
	"time"
	"users_service/pkg/logger"

	pb "users_service/genproto/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewAuthRepo(db *pgxpool.Pool, log logger.ILogger) *authRepo {
	return &authRepo{
		db:  db,
		log: log,
	}
}

func (a *authRepo) Create(ctx context.Context, request *pb.CreateUser) (*pb.User, error) {

	var (
		user      = pb.User{}
		query     string
		err       error
		timeNow   = time.Now()
		createdAt time.Time
	)

	query = `insert into users (
		email,
		password_hash,
		full_name,
		created_at
	) values ($1, $2, $3, $4) returning 
		id,
		email,
		full_name,
		user_role,
		created_at
	`

	if err = a.db.QueryRow(ctx, query,
		request.Email,
		request.Password,
		request.FullName,
		timeNow).
		Scan(
			&user.Id,
			&user.Email,
			&user.FullName,
			&user.UserRole,
			&createdAt,
		); err != nil {
		a.log.Error("error while creating user in storage layer", logger.Error(err))
		return nil, err
	}

	user.CreatedAt = createdAt.Format(Layout)

	return &user, nil
}

func (a *authRepo) GetByEmail(ctx context.Context, request *pb.Email) (*pb.UserByEmail, error) {

	var (
		user      = pb.UserByEmail{}
		query     string
		err       error
		createdAt time.Time
	)

	query = `
	select
		id,
		email,
		full_name,
		password_hash,
		user_role,
		created_at
	from 
		users 
	where
		email = $1 and
		deleted_at is null
	`

	if err = a.db.QueryRow(ctx, query, request.GetEmail()).Scan(
		&user.Id,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.UserRole,
		&createdAt,
	); err != nil {
		a.log.Error("error while getting user id by username", logger.Error(err))
		return nil, err
	}

	user.CreatedAt = createdAt.Format(Layout)

	return &user, nil
}

func (a *authRepo) DeleteRefreshTokenByUserId(ctx context.Context, request *pb.PrimaryKey) (*pb.Void, error) {

	var (
		query string
		err   error
	)

	query = `
		delete from 
			refresh_tokens
		where
			user_id = $1
	`

	if _, err = a.db.Exec(ctx, query, request.Id); err != nil {
		a.log.Error("error while deleting user's refresh token from toble", logger.Error(err))
		return &pb.Void{}, err
	}
	return &pb.Void{}, nil
}

func (a *authRepo) StoreRefreshToken(ctx context.Context, request *pb.RefreshToken) (*pb.Void, error) {

	var (
		query string
		err   error
		// expiresIn time.Time
	)

	query = `
	insert into refresh_tokens (
		user_id,
		refresh_token,
		expires_in
	) values ($1, $2, $3)
	`
	// expiresIn, err = time.Parse(time.RFC3339, request.ExpiresIn)
	// if err != nil {
	// 	return &pb.Void{}, err
	// }

	if _, err = a.db.Exec(ctx, query,
		request.UserId,
		request.RefreshToken,
		request.ExpiresIn,
	); err != nil {
		return &pb.Void{}, err
	}

	return &pb.Void{}, nil
}

func (a *authRepo) CheckRefreshTokenExists(ctx context.Context, request *pb.RequestRefreshToken) (*pb.Void, error) {

	var (
		query string
		err   error
		exist int
	)

	query = `
		select
			1
		from
			refresh_tokens
		where
			refresh_token = $1
	`

	err = a.db.QueryRow(ctx, query, request.RefreshToken).Scan(&exist)

	if err != nil && err.Error() != "no rows in result set" {
		a.log.Error("error while checking refresh token is exists", logger.Error(err))
		return &pb.Void{}, err
	}

	if exist != 1 {
		a.log.Error("error: refresh token not found in database")
		return &pb.Void{}, fmt.Errorf("error: refresh token not found in database")
	}

	return &pb.Void{}, nil
}

func (a *authRepo) CheckEmailExists(ctx context.Context, request *pb.Email) (*pb.Void, error) {

	var exist int

	query := `
		select
			1
		from
			users
		where
			email = $1
	`

	err := a.db.QueryRow(ctx, query, request.Email).Scan(&exist)
	if err != nil && err.Error() != "no rows in result set" {
		a.log.Error("error while checking if email exists", logger.Error(err))
		return &pb.Void{}, err
	}

	if exist != 1 {
		a.log.Error("error: email not found in database")
		return &pb.Void{}, fmt.Errorf("email not found in database")
	}

	return &pb.Void{}, nil
}

func (a *authRepo) ResetPassword(ctx context.Context, request *pb.ResetPassword) (*pb.Void, error) {

	var (
		query string
		err   error
	)

	query = `
		update 
			users
		set
			password_hash = $1
		where
			email = $2
	`

	if _, err = a.db.Exec(ctx, query,
		request.NewPassword,
		request.Email,
	); err != nil {
		a.log.Error("error while saving new password in storage layer", logger.Error(err))
		return &pb.Void{}, err
	}

	return &pb.Void{}, nil
}
