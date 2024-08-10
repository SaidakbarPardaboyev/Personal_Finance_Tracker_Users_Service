package postgres

import (
	"auth_service/configs"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, cfg configs.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresName,
	)

	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	conn.MaxConns = 100

	db, err := pgxpool.NewWithConfig(ctx, conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
