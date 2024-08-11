package redis

import (
	"users_service/pkg/logger"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type usersRedisRepo struct {
	client redis.Client
	log    logger.ILogger
}

func NewUsersRedisRepo(client *redis.Client, log logger.ILogger) *usersRedisRepo {
	return &usersRedisRepo{
		client: *client,
		log:    log,
	}
}

func (u *usersRedisRepo) SaveCodeWithEmail(ctx context.Context, email, code string) error {

	var (
		err error
	)

	if err = u.client.Set(ctx, email, code, time.Minute*2).Err(); err != nil {
		return err
	}

	return nil
}

func (u *usersRedisRepo) GetCodeWithEmail(ctx context.Context, email string) (string, error) {

	var (
		code string
		err  error
	)

	code, err = u.client.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}

	return code, nil
}
