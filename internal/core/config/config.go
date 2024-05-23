package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	AppDeviceHeaderKey     string = "x-AppDevice"
	AppVersionHeaderKey    string = "x-AppVersion"
	AuthorizationHeaderKey string = "Authorization"
	AuthUserUUIDKey        string = "UserUUID"
)

var (
	once   sync.Once
	config Config
)

type App struct {
	Name               string
	Env                string
	Version            string
	URL                string
	Port               string
	Debug              bool
	Timezone           string
	Locale             string
	FallbackLocale     string
	PathLocale         string
	GracefullyShutdown time.Duration
}

type Log struct {
	FilePath   string
	Level      string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type SwaggerInfo struct {
	Title       string
	Description string
	Version     string
}

type Swagger struct {
	Host     string
	Schemes  string
	Info     SwaggerInfo
	Enable   bool
	Username string
	Password string
}

type DBPostgres struct {
	SSLMode            string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxLifetime        time.Duration
	Timezone           string
}

type DB struct {
	Connection string
	Host       string
	Port       string
	Name       string
	Username   string
	Password   string
	Postgres   DBPostgres
}

type Redis struct {
	Host               string
	Port               string
	Password           string
	DB                 int
	Prefix             string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

type Jwt struct {
	AccessTokenSecret    string
	AccessTokenExpireDay time.Duration
}

// Config represents the application configuration.
type Config struct {
	App     App
	Log     Log
	Swagger Swagger
	DB      DB
	Redis   Redis
	Jwt     Jwt
}

// LoadConfig loads configuration from .env file and populates the Config struct.
func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %v", err)
	}

	var app App
	app.Name = os.Getenv("APP_NAME")
	app.Env = os.Getenv("APP_ENV")
	app.Version = os.Getenv("APP_VERSION")
	app.URL = os.Getenv("APP_URL")
	app.Port = os.Getenv("APP_PORT")
	app.Debug = getBoolEnv("APP_DEBUG")
	app.Timezone = os.Getenv("APP_TIMEZONE")
	app.Locale = os.Getenv("APP_LOCALE")
	app.FallbackLocale = os.Getenv("APP_FALLBACK_LOCALE")
	app.PathLocale = os.Getenv("APP_PATH_LOCALE")
	app.GracefullyShutdown = time.Duration(getIntEnv("APP_GRACEFULLY_SHUTDOWN", 5))

	var log Log
	log.FilePath = os.Getenv("LOG_FILE_PATH")
	log.Level = os.Getenv("LOG_LEVEL")
	log.MaxSize = getIntEnv(os.Getenv("LOG_MAX_SIZE"), 1)
	log.MaxAge = getIntEnv(os.Getenv("LOG_MAX_AGE"), 5)
	log.MaxBackups = getIntEnv(os.Getenv("LOG_MAX_BACKUPS"), 10)

	var swagger Swagger
	swagger.Host = os.Getenv("SWAGGER_HOST")
	swagger.Schemes = os.Getenv("SWAGGER_SCHEMES")
	swagger.Info.Title = os.Getenv("SWAGGER_INFO_TITLE")
	swagger.Info.Description = os.Getenv("SWAGGER_INFO_DESCRIPTION")
	swagger.Info.Version = os.Getenv("SWAGGER_INFO_VERSION")
	swagger.Enable = getBoolEnv("SWAGGER_ENABLE")
	swagger.Username = os.Getenv("SWAGGER_USERNAME")
	swagger.Password = os.Getenv("SWAGGER_PASSWORD")

	var db DB
	db.Connection = os.Getenv("DB_CONNECTION")
	db.Host = os.Getenv("DB_HOST")
	db.Port = os.Getenv("DB_PORT")
	db.Name = os.Getenv("DB_NAME")
	db.Username = os.Getenv("DB_USERNAME")
	db.Password = os.Getenv("DB_PASSWORD")
	db.Postgres.SSLMode = os.Getenv("DB_POSTGRES_SSL_MODE")
	db.Postgres.MaxOpenConnections = getIntEnv("DB_POSTGRES_MAX_OPEN_CONNECTIONS", 0)
	db.Postgres.MaxIdleConnections = getIntEnv("DB_POSTGRES_MAX_IDLE_CONNECTIONS", 0)
	db.Postgres.MaxLifetime = time.Duration(getIntEnv("DB_POSTGRES_MAX_LIFETIME", 0))
	db.Postgres.Timezone = os.Getenv("DB_POSTGRES_TIMEZONE")

	var redis Redis
	redis.Host = os.Getenv("REDIS_HOST")
	redis.Port = os.Getenv("REDIS_PORT")
	redis.Password = os.Getenv("REDIS_PASSWORD")
	redis.DB = getIntEnv("REDIS_DB", 0)
	redis.Prefix = os.Getenv("REDIS_PREFIX")
	redis.DialTimeout = time.Duration(getIntEnv("REDIS_DIAL_TIMEOUT", 0))
	redis.ReadTimeout = time.Duration(getIntEnv("REDIS_READ_TIMEOUT", 0))
	redis.WriteTimeout = time.Duration(getIntEnv("REDIS_WRITE_TIMEOUT", 0))
	redis.PoolSize = getIntEnv("REDIS_POOL_SIZE", 0)
	redis.PoolTimeout = time.Duration(getIntEnv("REDIS_POOL_TIMEOUT", 0))
	redis.IdleTimeout = time.Duration(getIntEnv("REDIS_IDLE_TIMEOUT", 0))
	redis.IdleCheckFrequency = time.Duration(getIntEnv("REDIS_IDLE_CHECK_FREQUENCY", 0))

	var jwt Jwt
	jwt.AccessTokenSecret = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	jwt.AccessTokenExpireDay = time.Duration(getIntEnv("JWT_ACCESS_TOKEN_EXPIRE_DAY", 7))

	return Config{
		App:     app,
		Log:     log,
		Swagger: swagger,
		DB:      db,
		Redis:   redis,
		Jwt:     jwt,
	}, nil
}

// Helper function to convert string environment variable to bool
func getBoolEnv(key string) bool {
	val, _ := strconv.ParseBool(os.Getenv(key))
	return val
}

// Helper function to convert string environment variable to int
func getIntEnv(key string, defaultValue int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return val
}

func GetConfig() Config {
	once.Do(func() {
		var err error
		config, err = LoadConfig()
		if err != nil {
			panic(err)
		}
	})
	return config
}
