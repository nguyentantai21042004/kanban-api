package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/minio"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/postgre"
	"gitlab.com/tantai-kanban/kanban-api/internal/httpserver"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	pkgCrt "gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

// @title Kanban API
// @description This is the API documentation for Kanban.
// @version 1
// @host https://kanban-api.ngtantai.pro
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
	// MESSAGE QUEUE CONFIGURATION
	// =============================================================================

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

		// Storage Configuration
		MinIOClient: minioClient,

		// Authentication & Security Configuration
		JwtSecretKey: cfg.JWT.SecretKey,
		Encrypter:    encrypter,
		InternalKey:  cfg.InternalConfig.InternalKey,

		// WebSocket Configuration
		WebSocketConfig: cfg.WebSocket,

		// Monitoring & Notification Configuration
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
		log.Println("Shutting down gracefully...")

		log.Println("Cleanup completed")
		os.Exit(0)
	}()
}
