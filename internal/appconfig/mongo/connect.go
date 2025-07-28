package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.com/tantai-kanban/kanban-api/config"
	"gitlab.com/tantai-kanban/kanban-api/pkg/mongo"
)

const (
	connectTimeout = 10 * time.Second
)

// Connect connects to the database
func Connect(mongoConfig config.MongoConfig, uri string) (mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), connectTimeout)
	defer cancelFunc()

	opts := mongo.NewClientOptions().
		ApplyURI(uri)

	if mongoConfig.EnableMonitor {
		opts.SetMonitor(mongo.CommandMonitor{
			Started: func(ctx context.Context, e *mongo.CommandStartedEvent) {
			},
			Succeeded: func(ctx context.Context, e *mongo.CommandSucceededEvent) {
			},
			Failure: func(ctx context.Context, e *mongo.CommandFailedEvent) {
			},
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	err = client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to DB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

// Disconnect disconnects from the database.
func Disconnect(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
