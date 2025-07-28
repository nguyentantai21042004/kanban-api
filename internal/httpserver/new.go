package httpserver

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/oauth"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgCrt "gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/minio"
	"gitlab.com/tantai-kanban/kanban-api/pkg/mongo"
	pkgRabbitMQ "gitlab.com/tantai-kanban/kanban-api/pkg/rabbitmq"
	"gitlab.com/tantai-kanban/kanban-api/pkg/redis"
)

type HTTPServer struct {
	// Server Configuration
	gin  *gin.Engine
	l    pkgLog.Logger
	host string
	port int
	mode string

	// Database Configuration
	postgresDB *sql.DB
	mongoDB    mongo.Client

	// Cache Configuration
	redisClient *redis.Client

	// Message Queue Configuration
	amqpConn *pkgRabbitMQ.Connection

	// Storage Configuration
	minioClient minio.MinIO

	// Authentication & Security Configuration
	jwtSecretKey string
	encrypter    pkgCrt.Encrypter
	internalKey  string
	oauthConfig  oauth.OauthConfig

	// External Services Configuration
	smtpConfig config.SMTPConfig

	// Monitoring & Notification Configuration
	telegram TeleCredentials
	discord  *discord.DiscordWebhook
}

type Config struct {
	// Server Configuration
	Logger pkgLog.Logger
	Host   string
	Port   int
	Mode   string

	// Database Configuration
	PostgresDB *sql.DB
	MongoDB    mongo.Client

	// Cache Configuration
	RedisClient *redis.Client

	// Message Queue Configuration
	AMQPConn *pkgRabbitMQ.Connection

	// Storage Configuration
	MinIOClient minio.MinIO

	// Authentication & Security Configuration
	JwtSecretKey string
	Encrypter    pkgCrt.Encrypter
	InternalKey  string
	OauthConfig  oauth.OauthConfig

	// External Services Configuration
	SMTPConfig config.SMTPConfig

	// Monitoring & Notification Configuration
	Telegram      TeleCredentials
	DiscordConfig *discord.DiscordWebhook
}

type TeleCredentials struct {
	BotKey string
	ChatIDs
}

type ChatIDs struct {
	ReportBug int64
}

func New(l pkgLog.Logger, cfg Config) (*HTTPServer, error) {
	if cfg.Mode == productionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	h := &HTTPServer{
		// Server Configuration
		l:    l,
		gin:  gin.Default(),
		host: cfg.Host,
		port: cfg.Port,
		mode: cfg.Mode,

		// Database Configuration
		postgresDB: cfg.PostgresDB,
		mongoDB:    cfg.MongoDB,

		// Cache Configuration
		redisClient: cfg.RedisClient,

		// Message Queue Configuration
		amqpConn: cfg.AMQPConn,

		// Storage Configuration
		minioClient: cfg.MinIOClient,

		// Authentication & Security Configuration
		jwtSecretKey: cfg.JwtSecretKey,
		encrypter:    cfg.Encrypter,
		internalKey:  cfg.InternalKey,
		oauthConfig:  cfg.OauthConfig,

		// External Services Configuration
		smtpConfig: cfg.SMTPConfig,

		// Monitoring & Notification Configuration
		telegram: cfg.Telegram,
		discord:  cfg.DiscordConfig,
	}

	if err := h.validate(); err != nil {
		return nil, err
	}

	return h, nil
}

func (s HTTPServer) validate() error {
	requiredDeps := []struct {
		dep interface{}
		msg string
	}{
		// Server Configuration
		{s.l, "logger is required"},
		{s.mode, "mode is required"},
		{s.host, "host is required"},
		{s.port, "port is required"},

		// Database Configuration
		{s.postgresDB, "postgresDB is required"},

		// Cache Configuration
		{s.redisClient, "redisClient is required"},

		// Message Queue Configuration
		{s.amqpConn, "amqpConn is required"},

		// Storage Configuration
		{s.minioClient, "minioClient is required"},

		// Authentication & Security Configuration
		{s.jwtSecretKey, "jwtSecretKey is required"},
		{s.encrypter, "encrypter is required"},
		{s.internalKey, "internalKey is required"},
		{s.oauthConfig, "oauthConfig is required"},

		// External Services Configuration
		{s.smtpConfig, "smtpConfig is required"},

		// Monitoring & Notification Configuration
		{s.telegram, "telegram is required"},
		{s.discord, "discord is required"},
	}

	for _, dep := range requiredDeps {
		if dep.dep == nil {
			return errors.New(dep.msg)
		}
	}

	return nil
}
