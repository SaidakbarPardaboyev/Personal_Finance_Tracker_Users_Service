package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	ApiGatewayHttpHost string
	ApiGatewayHttpPort string

	UserServiceGrpcHost string
	UserServiceGrpcPort string

	// LearingServiceGrpcHost string
	// LearingServiceGrpcPort string

	// ProgresServiceGrpcHost string
	// ProgresServiceGrpcPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresName     string
	PostgresPassword string

	MongoDBHost     string
	MongoDBPort     string
	MongoDBUser     string
	MongoDBName     string
	MongoDBPassword string

	RedisHost     string
	RedisDBNumber int
	RedisPort     string
	RedisPassword string

	SigningKeyAccess  string
	SigningKeyRefresh string

	ServiceName string
	LoggerLevel string
	LogPath     string

	Email    string
	Password string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file not found: %s", err)
	}
	config := Config{}
	config.ApiGatewayHttpHost = cast.ToString(coalesce("API_GATEWAY_HTTP_HOST", "localhost"))
	config.ApiGatewayHttpPort = cast.ToString(coalesce("API_GATEWAY_HTTP_PORT", ":8080"))

	config.UserServiceGrpcHost = cast.ToString(coalesce("USER_SERVICE_GRPC_HOST", "localhost"))
	config.UserServiceGrpcPort = cast.ToString(coalesce("USER_SERVICE_GRPC_PORT", ":1111"))
	// config.UserServiceHttpHost = cast.ToString(coalesce("USER_SERVICE_HTTP_HOST", "localhost"))
	// config.UserServiceHttpPort = cast.ToString(coalesce("USER_SERVICE_HTTP_PORT", ":2222"))

	// config.LearingServiceGrpcHost = cast.ToString(coalesce("LEARNING_SERVICE_GRPC_HOST", "localhost"))
	// config.LearingServiceGrpcPort = cast.ToString(coalesce("LEARNING_SERVICE_GRPC_PORT", ":3333"))

	// config.ProgresServiceGrpcHost = cast.ToString(coalesce("PROGRESS_SERVICE_GRPC_HOST", "localhost"))
	// config.ProgresServiceGrpcPort = cast.ToString(coalesce("PROGRESS_SERVICE_GRPC_PORT", ":4444"))

	config.PostgresHost = cast.ToString(coalesce("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToString(coalesce("POSTGRES_PORT", "5432"))
	config.PostgresUser = cast.ToString(coalesce("POSTGRES_USER", "postgres"))
	config.PostgresName = cast.ToString(coalesce("POSTGRES_NAME", "language_leap_auth_service"))
	config.PostgresPassword = cast.ToString(coalesce("POSTGRES_PASSWORD", "root"))

	config.MongoDBHost = cast.ToString(coalesce("MONGODB_HOST", "localhost"))
	config.MongoDBPort = cast.ToString(coalesce("MONGODB_PORT", "5432"))
	config.MongoDBUser = cast.ToString(coalesce("MONGODB_USER", "mongo"))
	config.MongoDBName = cast.ToString(coalesce("MONGODB_NAME", "market_product_service"))
	config.MongoDBPassword = cast.ToString(coalesce("MONGODB_PASSWORD", ""))

	config.RedisHost = cast.ToString(coalesce("REDIS_HOST", "localhost"))
	config.RedisDBNumber = cast.ToInt(coalesce("REDIS_DBNUMBER", 0))
	config.RedisPort = cast.ToString(coalesce("REDIS_PORT", "6379"))
	config.RedisPassword = cast.ToString(coalesce("REDIS_PASSWORD", ""))

	config.SigningKeyAccess = cast.ToString(coalesce("SINGNING_KEY_ACCESS", "SSECCA"))
	config.SigningKeyRefresh = cast.ToString(coalesce("SINGNING_KEY_REFRESH", "HSERFER"))

	config.ServiceName = cast.ToString(coalesce("SERVICE_NAME", "auth_service"))
	config.LoggerLevel = cast.ToString(coalesce("LOGGER_LEVEL", "debug"))
	config.LogPath = cast.ToString(coalesce("LOG_PATH", "app.log"))

	config.Email = cast.ToString(coalesce("EMAIL", "s@gmail.com"))
	config.Password = cast.ToString(coalesce("PASSWORD", "nothing"))

	return &config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	if res, exists := os.LookupEnv(key); exists {
		return res
	}
	return defaultValue
}
