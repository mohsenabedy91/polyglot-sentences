package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"time"
)

var dbClient *sql.DB

func InitClient(ctx context.Context, log logger.Logger, cfg config.Config) error {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password,
		cfg.DB.Name, cfg.DB.Postgres.SSLMode, cfg.DB.Postgres.Timezone)

	if dbClient, err = sql.Open("postgres", dsn); err != nil {
		log.Error(logger.Database, logger.Startup, fmt.Sprintf("There is an Error in Open DB : %v", err), nil)
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = dbClient.PingContext(ctx); err != nil {
		log.Error(logger.Database, logger.Startup, fmt.Sprintf("Database Ping is not available %v", err), nil)
		return err
	}

	dbClient.SetMaxOpenConns(cfg.DB.Postgres.MaxOpenConnections)
	dbClient.SetMaxIdleConns(cfg.DB.Postgres.MaxIdleConnections)
	dbClient.SetConnMaxLifetime(cfg.DB.Postgres.MaxLifetime * time.Minute)

	log.Info(logger.Database, logger.Startup, "Database client initialized", nil)

	return nil
}

func Get() *sql.DB {
	return dbClient
}

func Close() error {
	if cErr := dbClient.Close(); cErr != nil {
		return cErr
	}

	return nil
}

func getMigrateInstance(ctx context.Context, log logger.Logger) (*migrate.Migrate, error) {
	// Get the application configuration
	cfg := config.GetConfig()

	// Initialize the database client
	if err := InitClient(ctx, log, cfg); err != nil {
		return nil, err
	}

	// Get the database client return *sql.DB
	db := Get()

	// Create a new postgres driver
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	// Create a new migrate instance
	return migrate.NewWithDatabaseInstance(
		"file://internal/adapter/storage/postgres/migrations",
		cfg.DB.Name,
		dbDriver,
	)
}

func RunMigrations(log logger.Logger) error {
	// Get the migration instance
	instance, err := getMigrateInstance(context.Background(), log)
	if err != nil {
		return err
	}

	// Run all up migrations
	if err = instance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func RunDownMigration(log logger.Logger, step int) error {
	// Get the migration instance
	instance, err := getMigrateInstance(context.Background(), log)
	if err != nil {
		return err
	}
	if step > 0 {
		step = -step
	}
	// Run down migration to revert the last migration
	if err = instance.Steps(step); err != nil {
		return err
	}

	return nil
}
