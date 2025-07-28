package main

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/internal/appconfig/postgre"
	"gitlab.com/tantai-kanban/kanban-api/internal/consumer"
	pkgCrt "gitlab.com/tantai-kanban/kanban-api/pkg/encrypter"
	pkgLog "gitlab.com/tantai-kanban/kanban-api/pkg/log"
	"gitlab.com/tantai-kanban/kanban-api/pkg/rabbitmq"
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

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	srv, err := consumer.New(l, consumer.ConsumerConfig{
		Encrypter:    crp,
		JwtSecretKey: cfg.JWT.SecretKey,
		InternalKey:  cfg.InternalConfig.InternalKey,
		PostgresDB:   postgresDB,
		AMQPConn:     conn,
	})
	if err != nil {
		panic(err)
	}

	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
