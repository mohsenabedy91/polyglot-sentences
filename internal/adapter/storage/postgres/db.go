package postgres

import (
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

func InitClient(config config.Config, log logger.Logger) error {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password,
		config.DB.Name, config.DB.Postgres.SSLMode, config.DB.Postgres.Timezone)

	if dbClient, err = sql.Open("postgres", dsn); err != nil {
		log.Error(logger.Database, logger.Startup, fmt.Sprintf("There is an Error in Open DB : %v", err), nil)
		return err
	}

	if err = dbClient.Ping(); err != nil {
		log.Error(logger.Database, logger.Startup, fmt.Sprintf("Database Ping is not available %v", err), nil)
		return err
	}

	dbClient.SetMaxOpenConns(config.DB.Postgres.MaxOpenConnections)
	dbClient.SetMaxIdleConns(config.DB.Postgres.MaxIdleConnections)
	dbClient.SetConnMaxLifetime(config.DB.Postgres.MaxLifetime * time.Minute)

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

func getMigrateInstance(log logger.Logger) (*migrate.Migrate, error) {
	// Get the application configuration
	cfg := config.GetConfig()

	// Initialize the database client
	if err := InitClient(cfg, log); err != nil {
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
	instance, err := getMigrateInstance(log)
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
	instance, err := getMigrateInstance(log)
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
