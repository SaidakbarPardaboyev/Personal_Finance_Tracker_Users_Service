package postgres

import (
	"context"
	"database/sql"
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
		password_hash,
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
			&user.Password,
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

func (a *authRepo) GetByEmail(ctx context.Context, request *pb.Email) (*pb.User, error) {

	var (
		user      = pb.User{}
		query     string
		err       error
		createdAt time.Time
	)

	query = `
	select
		id,
		email,
		password_hash,
		full_name,
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
		&user.Password,
		&user.FullName,
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
	)

	query = `
	insert into refresh_tokens (
		user_id,
		refresh_token,
		expires_in
	) values ($1, $2, $3)
	`

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
		exist = sql.NullInt64{}
	)

	query = `
		select
			1
		from
			refresh_tokens
		where
			refresh_token = $1
	`

	if err = a.db.QueryRow(ctx, query, request.RefreshToken).Scan(&exist); err != nil {
		a.log.Error("error user not found in users table", logger.Error(err))
		return &pb.Void{}, err
	}

	if !exist.Valid || exist.Int64 != 1 {
		a.log.Error("error user not found in users table")
		return &pb.Void{}, fmt.Errorf("error user not found in users table")
	}

	return &pb.Void{}, nil
}

func (a *authRepo) CheckEmailExists(ctx context.Context, request *pb.Email) (*pb.Void, error) {

	var (
		query string
		err   error
		exist = sql.NullInt64{}
	)

	query = `
		select
			1
		from
			users
		where
			email = $1
	`

	if err = a.db.QueryRow(ctx, query, request.Email).Scan(&exist); err != nil {
		a.log.Error("error user not found in users table", logger.Error(err))
		return &pb.Void{}, err
	}

	if !exist.Valid || exist.Int64 != 1 {
		a.log.Error("error user not found in users table")
		return &pb.Void{}, fmt.Errorf("error user not found in users table")
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
