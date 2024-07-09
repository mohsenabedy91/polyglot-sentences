//go:build !test

package setup

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

func InitializeDatabase(ctx context.Context, log logger.Logger, conf config.Config) (*sql.DB, error) {
	if err := postgres.InitClient(ctx, log, conf); err != nil {
		log.Fatal(logger.Database, logger.Startup, fmt.Sprintf("Failed to setup postgres, error: %v", err), nil)
		return nil, err
	}
	return postgres.Get(), nil
}

func InitializeQueue(log logger.Logger, conf config.Config) (*messagebroker.Queue, error) {
	queue := messagebroker.NewQueue(log, conf)

	driver, err := messagebroker.NewRabbitMQ(conf.RabbitMQ.URL, log)
	if err != nil {
		log.Fatal(logger.Queue, logger.Startup, fmt.Sprintf("Failed to setup queue, error: %v", err), nil)
		return nil, err
	}

	queue.Driver = driver

	return queue, nil
}
