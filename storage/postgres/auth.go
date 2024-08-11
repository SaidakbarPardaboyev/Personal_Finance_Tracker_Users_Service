package postgres

import (
	"users_service/models"
	"users_service/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (a *authRepo) Register(ctx context.Context, request *models.RequestRegister) (*models.ResponseRegister, error) {

	var (
		user    = models.ResponseRegister{}
		query   string
		err     error
		timeNow = time.Now()
	)

	query = `insert into users (
		username,
		email,
		password_hash,
		full_name,
		native_language,
		created_at
	) values ($1, $2, $3, $4, $5, $6) returning 
		id,
		username,
		email,
		full_name,
		native_language,
		created_at::text
	`

	if err = a.db.QueryRow(ctx, query,
		request.Username,
		request.Email,
		request.Password,
		request.FullName,
		request.NativeLanguage,
		timeNow).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.FullName,
			&user.NativeLanguage,
			&user.CreatedAt,
		); err != nil {
		a.log.Error("error while creating user in storage layer", logger.Error(err))
		return nil, err
	}

	return &user, nil
}

func (a *authRepo) GetUserByUsername(ctx context.Context, username string) (*models.UserForLogin, error) {

	var (
		user  = models.UserForLogin{}
		query string
		err   error
	)

	query = `
	select
		id,
		username,
		email,
		password_hash,
		full_name,
		created_at::text
	from 
		users 
	where
		username = $1
	`

	if err = a.db.QueryRow(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.CreatedAt,
	); err != nil {
		a.log.Error("error while getting user id by username", logger.Error(err))
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) DeleteRefreshTokenByUserId(ctx context.Context, userId string) error {

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

	if _, err = a.db.Exec(ctx, query, userId); err != nil {
		a.log.Error("error while deleting user's refresh token from toble", logger.Error(err))
		return err
	}
	return nil
}

func (a *authRepo) StoreRefreshToken(ctx context.Context, request *models.StoreRefreshToken) error {

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
		return err
	}

	return nil
}

func (a *authRepo) CheckRefreshTokenExists(ctx context.Context, refreshToken string) error {

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

	if err = a.db.QueryRow(ctx, query, refreshToken).Scan(&exist); err != nil {
		a.log.Error("error user not found in users table", logger.Error(err))
		return err
	}

	if !exist.Valid || exist.Int64 != 1 {
		a.log.Error("error user not found in users table")
		return fmt.Errorf("error user not found in users table")
	}

	return err
}

func (a *authRepo) CheckEmailExists(ctx context.Context, email string) error {

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

	if err = a.db.QueryRow(ctx, query, email).Scan(&exist); err != nil {
		a.log.Error("error user not found in users table", logger.Error(err))
		return err
	}

	if !exist.Valid || exist.Int64 != 1 {
		a.log.Error("error user not found in users table")
		return fmt.Errorf("error user not found in users table")
	}

	return err
}

func (a *authRepo) ResetPassword(ctx context.Context, email, password string) error {

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
		password,
		email,
	); err != nil {
		a.log.Error("error while saving new password in storage layer", logger.Error(err))
		return err
	}

	return nil
}
