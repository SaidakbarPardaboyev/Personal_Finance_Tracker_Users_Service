package storage

import (
	"context"
	"users_service/configs"
	"users_service/pkg/logger"
	"users_service/storage/postgres"
	rds "users_service/storage/redis"

	"github.com/redis/go-redis/v9"

	pb "users_service/genproto/users"

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
	Create(context.Context, *pb.CreateUser) (*pb.User, error)
	GetByEmail(context.Context, *pb.Email) (*pb.User, error)
	DeleteRefreshTokenByUserId(context.Context, *pb.PrimaryKey) (*pb.Void, error)
	StoreRefreshToken(context.Context, *pb.RefreshToken) (*pb.Void, error)
	CheckRefreshTokenExists(context.Context, *pb.RequestRefreshToken) (*pb.Void, error)
	CheckEmailExists(context.Context, *pb.Email) (*pb.Void, error)
	ResetPassword(context.Context, *pb.ResetPassword) (*pb.Void, error)
}

type IUsersStorage interface {
	GetById(context.Context, *pb.PrimaryKey) (*pb.User, error)
	GetAll(context.Context, *pb.GetListRequest) (*pb.Users, error)
	Update(context.Context, *pb.UpdateUser) (*pb.UpdatedUser, error)
	Delete(context.Context, *pb.PrimaryKey) (*pb.Void, error)
	CheckPasswordExisis(context.Context, *pb.ChangePassword) (bool, error)
	ChangePassword(context.Context, *pb.ChangePassword) (*pb.Void, error)
	ChangeUserRole(context.Context, *pb.ChangeUserRole) (*pb.Void, error)
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
