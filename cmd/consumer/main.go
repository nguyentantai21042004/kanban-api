package main

import (
	"context"

	_ "github.com/lib/pq"
	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/oauth"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/postgre"
	"gitlab.com/tantai-kanban/kanban-api/internal/consumer"
	pkgCrt "gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/rabbitmq"
	pkgRedis "gitlab.com/tantai-kanban/kanban-api/pkg/redis"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	crp := pkgCrt.NewEncrypter(cfg.Encrypter.Key)

	postgresDB, err := postgre.Connect(context.Background(), cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer postgre.Disconnect(context.Background(), postgresDB)

	conn, err := rabbitmq.Dial(cfg.RabbitMQConfig.URL, true)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	redisClient, err := pkgRedis.Connect(pkgRedis.NewClientOptions().SetOptions(cfg.Redis))
	if err != nil {
		panic(err)
	}
	defer redisClient.Disconnect()

	oauthConfig := oauth.NewOauthConfig(cfg.Oauth)

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	srv, err := consumer.New(l, consumer.ConsumerConfig{
		Encrypter:    crp,
		JwtSecretKey: cfg.JWT.SecretKey,
		Telegram: consumer.TeleCredentials{
			BotKey: cfg.Telegram.BotKey,
			ChatIDs: consumer.ChatIDs{
				ReportBug: cfg.Telegram.ChatIDs.ReportBug,
			},
		},
		InternalKey: cfg.InternalConfig.InternalKey,
		PostgresDB:  postgresDB,
		SMTPConfig:  cfg.SMTP,
		AMQPConn:    conn,
		RedisClient: &redisClient,
		OauthConfig: oauthConfig,
	})
	if err != nil {
		panic(err)
	}

	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
