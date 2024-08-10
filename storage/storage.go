package storage

import (
	"auth_service/configs"
	"auth_service/models"
	"auth_service/pkg/logger"
	"auth_service/storage/postgres"
	rds "auth_service/storage/redis"
	"context"

	"github.com/redis/go-redis/v9"

	pb "auth_service/genproto/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbPostgres  *pgxpool.Pool
	redisClient *redis.Client
	log         logger.ILogger
}

type IStorage interface {
	Close()
	Auth() IAuthStorage
	Users() IUsersStorage
	RedisClient() IUserRedisStorage
}

type IAuthStorage interface {
	Register(context.Context, *models.RequestRegister) (*models.ResponseRegister, error)
	GetUserByUsername(context.Context, string) (*models.UserForLogin, error)
	DeleteRefreshTokenByUserId(context.Context, string) error
	StoreRefreshToken(context.Context, *models.StoreRefreshToken) error
	CheckRefreshTokenExists(context.Context, string) error
	CheckEmailExists(context.Context, string) error
	ResetPassword(context.Context, string, string) error
}

type IUsersStorage interface {
	GetUserProfile(context.Context, *pb.PrimaryKey) (*pb.User, error)
	UpdateUserProfile(context.Context, *pb.UpdateUser) (*pb.UpdateProfileResponce, error)
	CheckPasswordExisis(context.Context, *pb.ChangePassword) (bool, error)
	ChangePassword(context.Context, *pb.ChangePassword) (*pb.Message, error)
}

type IUserRedisStorage interface {
	SaveCodeWithEmail(context.Context, string, string) error
	GetCodeWithEmail(context.Context, string) (string, error)
}

func New(ctx context.Context, cfg *configs.Config, log *logger.ILogger) (IStorage, error) {
	dbPostgres, err := postgres.ConnectDB(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := rds.NewRedisClient()
	if err != nil {
		return nil, err
	}

	return &Storage{
		dbPostgres:  dbPostgres,
		redisClient: redisClient,
		log:         *log,
	}, nil
}

func (s *Storage) Close() {
	s.dbPostgres.Close()
}

func (s *Storage) Auth() IAuthStorage {
	return postgres.NewAuthRepo(s.dbPostgres, s.log)
}

func (s *Storage) Users() IUsersStorage {
	return postgres.NewUsersRepo(s.dbPostgres, s.log)
}

func (s *Storage) RedisClient() IUserRedisStorage {
	return rds.NewUsersRedisRepo(s.redisClient, s.log)
}
