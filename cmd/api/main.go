package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/minio"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/mongo"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/postgre"
	"gitlab.com/tantai-kanban/kanban-api/internal/httpserver"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgCrt "gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/rabbitmq"
)

// @title SMAP Authenticate API
// @description This is the API documentation for SMAP Authenticate.
// @description authenticate
// @description `110001 ("Wrong query"),`
// @description `110002 ("Wrong body"),`
// @description `110003 ("User not found"),`
// @description `110004 ("Email existed"),`
// @description `110005 ("Wrong password"),`
// @version 1
// @host 192.168.1.216
// @schemes https
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown()

	// Initialize logger first (needed for all other services)
	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	// Initialize Encrypter
	encrypter := pkgCrt.NewEncrypter(cfg.Encrypter.Key)

	// =============================================================================
	// DATABASE CONFIGURATION
	// =============================================================================

	// Initialize PostgreSQL
	postgresDB, err := postgre.Connect(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	defer postgre.Disconnect(context.Background(), postgresDB)

	// =============================================================================
	// CACHE CONFIGURATION
	// =============================================================================

	uri, err := encrypter.Decrypt(cfg.Mongo.EncodedURI)
	if err != nil {
		log.Fatal("Failed to decrypt mongo uri: ", err)
	}
	mongoDB, err := mongo.Connect(cfg.Mongo, uri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer mongo.Disconnect(mongoDB)

	// Initialize Redis
	// redisOpts := redis.NewClientOptions().SetOptions(cfg.Redis)
	// redisClient, err := redis.Connect(redisOpts)
	// if err != nil {
	// 	log.Fatal("Failed to connect to Redis: ", err)
	// }
	// defer redisClient.Disconnect()
	// log.Println("✅ Redis connected successfully")

	// =============================================================================
	// MESSAGE QUEUE CONFIGURATION
	// =============================================================================

	// Initialize RabbitMQ
	// log.Println("Connecting to RabbitMQ...")
	amqpConn, err := rabbitmq.Dial(cfg.RabbitMQConfig.URL, true)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer amqpConn.Close()

	// =============================================================================
	// STORAGE CONFIGURATION
	// =============================================================================

	// Initialize MinIO
	// log.Println("Connecting to MinIO...")
	minioClient, err := minio.Connect(context.Background(), cfg.MinIO)
	if err != nil {
		log.Fatal("Failed to connect to MinIO: ", err)
	}
	defer minio.Close()

	// =============================================================================
	// AUTHENTICATION & SECURITY CONFIGURATION
	// =============================================================================

	// =============================================================================
	// EXTERNAL SERVICES CONFIGURATION
	// =============================================================================

	// SMTP is configured via config, no connection needed

	// =============================================================================
	// MONITORING & NOTIFICATION CONFIGURATION
	// =============================================================================

	// Initialize Discord Webhook
	discordWebhook, err := discord.NewDiscordWebhook(cfg.Discord.ReportBugID, cfg.Discord.ReportBugToken)
	if err != nil {
		log.Fatal("Failed to initialize Discord webhook: ", err)
	}

	// =============================================================================
	// HTTP SERVER CONFIGURATION
	// =============================================================================

	srv, err := httpserver.New(l, httpserver.Config{
		// Server Configuration
		Logger: l,
		Host:   cfg.HTTPServer.Host,
		Port:   cfg.HTTPServer.Port,
		Mode:   cfg.HTTPServer.Mode,

		// Database Configuration
		PostgresDB: postgresDB,
		MongoDB:    mongoDB,

		// Cache Configuration
		// RedisClient: redisClient.(*redis.Client),

		// Message Queue Configuration
		AMQPConn: amqpConn,

		// Storage Configuration
		MinIOClient: minioClient,

		// Authentication & Security Configuration
		JwtSecretKey: cfg.JWT.SecretKey,
		Encrypter:    encrypter,
		InternalKey:  cfg.InternalConfig.InternalKey,
		// OauthConfig:  cfg.Oauth, // Add OAuth config

		// External Services Configuration
		SMTPConfig: cfg.SMTP,

		// Monitoring & Notification Configuration
		Telegram: httpserver.TeleCredentials{
			BotKey: cfg.Telegram.BotKey,
			ChatIDs: httpserver.ChatIDs{
				ReportBug: cfg.Telegram.ChatIDs.ReportBug,
			},
		},
		DiscordConfig: discordWebhook,
	})
	if err != nil {
		log.Fatal("Failed to initialize HTTP server: ", err)
	}

	// =============================================================================
	// START SERVER
	// =============================================================================

	if err := srv.Run(); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("🛑 Shutting down gracefully...")

		log.Println("✅ Cleanup completed")
		os.Exit(0)
	}()
}
