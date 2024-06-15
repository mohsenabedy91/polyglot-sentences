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
)

// please never change this keys
const (
	AuthTokenUserUUID       string = "sub"
	AuthTokenJTI            string = "jti"
	AuthTokenIssuedAt       string = "iat"
	AuthTokenExpirationTime string = "exp"
)

var (
	once   sync.Once
	config Config
)

type Kong struct {
	APIBaseUrl      string
	SwaggerFilePath string
	AuthorizeUrl    string
}

type App struct {
	ResetPasswordURL   string
	VerificationURL    string
	SupportEmail       string
	Env                string
	Debug              bool
	Timezone           string
	Locale             string
	FallbackLocale     string
	PathLocale         string
	GracefullyShutdown time.Duration
}

type Auth struct {
	Name    string
	Version string
	URL     string
	Port    string
	Debug   bool
}

type UserManagement struct {
	Name     string
	Version  string
	URL      string
	HTTPPort string
	GRPCPort string
	Debug    bool
}

type Profile struct {
	Debug bool
	Port  string
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

type OTP struct {
	ExpireSecond               time.Duration
	ForgetPasswordExpireSecond time.Duration
	Digits                     int8
}
type RabbitMQ struct {
	URL string
}

type SendGrid struct {
	Key     string
	Name    string
	Address string
}

type Oauth struct {
	Google
}

type Google struct {
	ClientId     string
	ClientSecret string
	CallbackURL  string
}

// Config represents the application configuration.
type Config struct {
	Kong           Kong
	App            App
	Auth           Auth
	UserManagement UserManagement
	Profile        Profile
	Log            Log
	Swagger        Swagger
	DB             DB
	Redis          Redis
	Jwt            Jwt
	OTP            OTP
	RabbitMQ       RabbitMQ
	SendGrid       SendGrid
	Oauth          Oauth
}

// LoadConfig loads configuration from .env file and populates the Config struct.
func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %v", err)
	}

	var kong Kong
	kong.APIBaseUrl = os.Getenv("KONG_API_BASE_URL")
	kong.SwaggerFilePath = os.Getenv("KONG_SWAGGER_FILE_PATH")
	kong.AuthorizeUrl = os.Getenv("KONG_AUTHORIZE_URL")

	var app App
	app.ResetPasswordURL = os.Getenv("RESET_PASSWORD_URL")
	app.VerificationURL = os.Getenv("APP_VERIFICATION_URL")
	app.SupportEmail = os.Getenv("APP_SUPPORT_EMAIL")
	app.Env = os.Getenv("APP_ENV")
	app.Debug = getBoolEnv("APP_DEBUG")
	app.Timezone = os.Getenv("APP_TIMEZONE")
	app.Locale = os.Getenv("APP_LOCALE")
	app.FallbackLocale = os.Getenv("APP_FALLBACK_LOCALE")
	app.PathLocale = os.Getenv("APP_PATH_LOCALE")
	app.GracefullyShutdown = time.Duration(getIntEnv("APP_GRACEFULLY_SHUTDOWN", 5))

	var auth Auth
	auth.Name = os.Getenv("AUTH_NAME")
	auth.Version = os.Getenv("AUTH_VERSION")
	auth.URL = os.Getenv("AUTH_URL")
	auth.Port = os.Getenv("AUTH_PORT")
	auth.Debug = getBoolEnv("AUTH_DEBUG")

	var userManagement UserManagement
	userManagement.Name = os.Getenv("USER_MANAGEMENT_NAME")
	userManagement.Version = os.Getenv("USER_MANAGEMENT_VERSION")
	userManagement.URL = os.Getenv("USER_MANAGEMENT_URL")
	userManagement.HTTPPort = os.Getenv("USER_MANAGEMENT_HTTP_PORT")
	userManagement.GRPCPort = os.Getenv("USER_MANAGEMENT_GRPC_PORT")
	userManagement.Debug = getBoolEnv("USER_MANAGEMENT_DEBUG")

	var profile Profile
	profile.Debug = getBoolEnv("PROFILE_DEBUG")
	profile.Port = os.Getenv("PROFILE_PORT")

	var log Log
	log.FilePath = os.Getenv("LOG_FILE_PATH")
	log.Level = os.Getenv("LOG_LEVEL")
	log.MaxSize = getIntEnv("LOG_MAX_SIZE", 1)
	log.MaxAge = getIntEnv("LOG_MAX_AGE", 5)
	log.MaxBackups = getIntEnv("LOG_MAX_BACKUPS", 10)

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

	var otp OTP
	otp.ExpireSecond = time.Duration(getIntEnv("OTP_EXPIRE_SECOND", 7)) * time.Second
	otp.ForgetPasswordExpireSecond = time.Duration(getIntEnv("FORGET_PASSWORD_EXPIRE_SECOND", 86400)) * time.Second
	otp.Digits = int8(getIntEnv("OTP_DIGITS", 6))

	var rabbitMQ RabbitMQ
	rabbitMQ.URL = os.Getenv("RABBITMQ_URL")

	var sendGrid SendGrid
	sendGrid.Key = os.Getenv("SEND_GRID_KEY")
	sendGrid.Name = os.Getenv("SEND_GRID_NAME")
	sendGrid.Address = os.Getenv("SEND_GRID_ADDRESS")

	var oauth Oauth
	oauth.Google.ClientId = os.Getenv("OAUTH_GOOGLE_CLIENT_ID")
	oauth.Google.ClientSecret = os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET")
	oauth.Google.CallbackURL = os.Getenv("OAUTH_GOOGLE_CALLBACK_URL")

	return Config{
		Kong:           kong,
		App:            app,
		Auth:           auth,
		Profile:        profile,
		Log:            log,
		Swagger:        swagger,
		DB:             db,
		Redis:          redis,
		Jwt:            jwt,
		UserManagement: userManagement,
		OTP:            otp,
		RabbitMQ:       rabbitMQ,
		SendGrid:       sendGrid,
		Oauth:          oauth,
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
