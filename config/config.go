package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	// Server Configuration
	HTTPServer HTTPServerConfig
	Logger     LoggerConfig

	// Database Configuration
	Postgres PostgresConfig
	Mongo    MongoConfig

	// Cache Configuration
	Redis RedisConfig

	// Message Queue Configuration
	RabbitMQConfig RabbitMQConfig

	// Storage Configuration
	MinIO MinIOConfig

	// Authentication & Security Configuration
	JWT            JWTConfig
	Encrypter      EncrypterConfig
	InternalConfig InternalConfig
	Oauth          OauthConfig
	GoogleDrive    GoogleDriveConfig

	// External Services Configuration
	SMTP SMTPConfig

	// Monitoring & Notification Configuration
	Telegram TelegramConfig
	Discord  DiscordConfig
}

// JWTConfig is the configuration for the JWT,
// which is used to generate and verify the JWT.
type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET"`
}

// HTTPServerConfig is the configuration for the HTTP server,
// which is used to start, call API, etc.
type HTTPServerConfig struct {
	Host string `env:"HOST" envDefault:""`
	Port int    `env:"APP_PORT" envDefault:"80"`
	Mode string `env:"API_MODE" envDefault:"debug"`
}

// LoggerConfig is the configuration for the logger,
// which is used to log the application.
type LoggerConfig struct {
	Level    string `env:"LOGGER_LEVEL" envDefault:"debug"`
	Mode     string `env:"LOGGER_MODE" envDefault:"debug"`
	Encoding string `env:"LOGGER_ENCODING" envDefault:"console"`
}

// PostgresConfig is the configuration for the Postgres,
// which is used to connect to the Postgres.
type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"POSTGRES_DB" envDefault:"postgres"`
}

// MongoConfig is the configuration for the Mongo,
// which is used to connect to the Mongo.
type MongoConfig struct {
	Database      string `env:"MONGODB_DATABASE"`
	EncodedURI    string `env:"MONGODB_ENCODED_URI"`
	EnableMonitor bool   `env:"MONGODB_ENABLE_MONITORING" envDefault:"true"`
}

type MinIOConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	Region    string `env:"MINIO_REGION" envDefault:"us-east-1"`
	Bucket    string `env:"MINIO_BUCKET"`
}

type CloudinaryConfig struct {
	CloudName string `env:"CLOUDINARY_CLOUD_NAME"`
	APIKey    string `env:"CLOUDINARY_API_KEY"`
	APISecret string `env:"CLOUDINARY_API_SECRET"`
}

// TelegramConfig is the configuration for the Telegram,
// which is used to send message to the Telegram.
type TelegramConfig struct {
	BotKey  string `env:"TELEGRAM_BOT_KEY" envDefault:""`
	ChatIDs TeleChatIDs
}

// TeleChatIDs is the configuration for the Telegram chat IDs,
// which is used to send message to the Telegram.
type TeleChatIDs struct {
	ReportBug int64 `env:"TELEGRAM_REPORT_BUG" envDefault:"-378952570"`
}

type DiscordConfig struct {
	ReportBugID    string `env:"DISCORD_REPORT_BUG_ID"`
	ReportBugToken string `env:"DISCORD_REPORT_BUG_TOKEN"`
}

// EncrypterConfig is the configuration for the encrypter,
// which is used to encrypt and decrypt the data.
type EncrypterConfig struct {
	Key string `env:"ENCRYPT_KEY"`
}

// InternalConfig is the configuration for the internal,
// which is used to check the internal request.
type InternalConfig struct {
	InternalKey string `env:"INTERNAL_KEY"`
}

// RabbitMQConfig is the configuration for the RabbitMQ,
// which is used to connect to the RabbitMQ.
type RabbitMQConfig struct {
	URL string `env:"RABBITMQ_URL"`
}

// SMTPConfig is the configuration for the SMTP,
// which is used to send email.
type SMTPConfig struct {
	Host     string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	Port     int    `env:"SMTP_PORT" envDefault:"587"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	From     string `env:"SMTP_FROM"`
	FromName string `env:"SMTP_FROM_NAME"`
}

// RedisConfig is the configuration for the Redis,
// which is used to connect to the Redis.
type RedisConfig struct {
	RedisAddr         []string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword     string   `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB           int      `env:"REDIS_DB" envDefault:"0"`
	RedisStandAlone   bool     `env:"REDIS_STANDALONE" envDefault:"true"`
	RedisPoolSize     int      `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RedisPoolTimeout  int      `env:"REDIS_POOL_TIMEOUT" envDefault:"10"`
	RedisMinIdleConns int      `env:"REDIS_MIN_IDLE_CONNS" envDefault:"10"`
}

type OauthConfig struct {
	Google   GoogleOauthConfig
	Facebook FacebookOauthConfig
	Gitlab   GitlabOauthConfig
}

type GoogleOauthConfig struct {
	ClientID     string   `env:"GOOGLE_OAUTH_CLIENT_ID"`
	ClientSecret string   `env:"GOOGLE_OAUTH_CLIENT_SECRET"`
	RedirectURL  string   `env:"GOOGLE_OAUTH_REDIRECT_URL"`
	Scopes       []string `env:"GOOGLE_OAUTH_SCOPES"`
	AuthURL      string   `env:"GOOGLE_OAUTH_AUTH_URL"`
	TokenURL     string   `env:"GOOGLE_OAUTH_TOKEN_URL"`
	UserInfoURL  string   `env:"GOOGLE_OAUTH_USER_INFO_URL"`
}

type GoogleDriveConfig struct {
	ClientID     string `env:"GOOGLE_DRIVE_CLIENT_ID"`
	ClientSecret string `env:"GOOGLE_DRIVE_CLIENT_SECRET"`
	RedirectURL  string `env:"GOOGLE_DRIVE_REDIRECT_URL"`
}

type FacebookOauthConfig struct {
	ClientID     string   `env:"FACEBOOK_OAUTH_CLIENT_ID"`
	ClientSecret string   `env:"FACEBOOK_OAUTH_CLIENT_SECRET"`
	RedirectURL  string   `env:"FACEBOOK_OAUTH_REDIRECT_URL"`
	Scopes       []string `env:"FACEBOOK_OAUTH_SCOPES"`
	AuthURL      string   `env:"FACEBOOK_OAUTH_AUTH_URL"`
	TokenURL     string   `env:"FACEBOOK_OAUTH_TOKEN_URL"`
	UserInfoURL  string   `env:"FACEBOOK_OAUTH_USER_INFO_URL"`
}

type GitlabOauthConfig struct {
	ClientID     string   `env:"GITLAB_OAUTH_CLIENT_ID"`
	ClientSecret string   `env:"GITLAB_OAUTH_CLIENT_SECRET"`
	RedirectURL  string   `env:"GITLAB_OAUTH_REDIRECT_URL"`
	Scopes       []string `env:"GITLAB_OAUTH_SCOPES"`
	AuthURL      string   `env:"GITLAB_OAUTH_AUTH_URL"`
	TokenURL     string   `env:"GITLAB_OAUTH_TOKEN_URL"`
	UserInfoURL  string   `env:"GITLAB_OAUTH_USER_INFO_URL"`
}

// Load is the function to load the configuration from the environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
