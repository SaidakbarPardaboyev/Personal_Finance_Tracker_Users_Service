package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	redisClient := *redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &redisClient, nil
}
