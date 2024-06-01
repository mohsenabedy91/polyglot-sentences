package main

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/server"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/repository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

func main() {
	cfg := config.GetConfig()
	log := logger.NewLogger(cfg.UserManagement.Name, cfg.Log)

	defer func() {
		err := postgres.Close()
		if err != nil {
			log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		}
	}()
	if err := postgres.InitClient(log, cfg); err != nil {
		return
	}
	postgresDB := postgres.Get()

	trans := translation.NewTranslation(cfg.App.Locale)
	trans.GetLocalizer(cfg.App.Locale)

	userRepo := repository.NewUserRepository(log, postgresDB)
	userService := userservice.New(log, userRepo)

	s := server.NewUserGRPCServer(cfg.UserManagement, userService, trans)
	grpcServer, err := s.StartUserGRPCServer()
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return
	}

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Server ...", nil)

	grpcServer.GracefulStop()
}
