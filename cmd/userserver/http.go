package main

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/routes"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/repository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authorizationservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
// @description "Bearer <your-jwt-token>"
func main() {
	cfg := config.GetConfig()

	log := logger.NewLogger(cfg.UserManagement.Name, cfg.Log)

	defer func() {
		err := postgres.Close()
		if err != nil {
			log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		}
	}()
	if err := postgres.InitClient(cfg, log); err != nil {
		return
	}
	postgresDB := postgres.Get()

	trans := translation.NewTranslation(cfg.App.Locale)
	trans.GetLocalizer(cfg.App.Locale)

	healthHandler := handler.NewHealthHandler(trans)

	// Init router
	router, err := routes.NewRouter(cfg, log, trans, *healthHandler)
	if err != nil {
		log.Error(logger.Internal, logger.Startup, fmt.Sprintf("There is an error when run http: %v", err), nil)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(log, postgresDB)
	userService := userservice.New(log, userRepo)
	userHandler := handler.NewUserHandler(userService)

	permissionRepo := repository.NewPermissionRepository(log, postgresDB)
	accessControlService := authorizationservice.New(log, permissionRepo, userService)

	router = router.NewUserRouter(accessControlService, *userHandler)

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", cfg.UserManagement.URL, cfg.UserManagement.HTTPPort)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router.Handler(),
	}
	log.Info(logger.Internal, logger.Startup, "Starting the HTTP server", map[logger.ExtraKey]interface{}{
		logger.ListeningAddress: server.Addr,
	})

	router.Serve(server)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Server ...", nil)

	timeout := cfg.App.GracefullyShutdown * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(logger.Internal, logger.Shutdown, fmt.Sprintf("server Shutdown: %v", err), nil)
	}

	select {
	case <-ctx.Done():
		log.Info(logger.Internal, logger.Shutdown, "timeout of 5 seconds.", nil)
	}
	log.Info(logger.Internal, logger.Shutdown, "Server exiting", nil)
}