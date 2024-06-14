package main

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/cmd/setup"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/server"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

func main() {
	cfg := config.GetConfig()
	log := logger.NewLogger(cfg.UserManagement.Name, cfg.Log)

	ctx := context.Background()
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		}
	}()
	postgresDB, err := setup.InitializeDatabase(ctx, log, cfg)
	if err != nil {
		return
	}
	uowFactory := func() repository.UnitOfWork {
		return repository.NewUnitOfWork(log, postgresDB)
	}

	trans := translation.NewTranslation(cfg.App.Locale)
	trans.GetLocalizer(cfg.App.Locale)

	userService := userservice.New(log)

	s := server.NewUserGRPCServer(cfg.UserManagement, userService, uowFactory)
	grpcServer, err := s.StartUserGRPCServer()
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return
	}

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Server ...", nil)

	grpcServer.GracefulStop()
}
