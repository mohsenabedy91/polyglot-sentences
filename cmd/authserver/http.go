package main

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/client"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/routes"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/repository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/aclservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/otpservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/oauth"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
// @description "Bearer <your-jwt-token>"
func main() {
	cfg := config.GetConfig()
	log := logger.NewLogger(cfg.Auth.Name, cfg.Log)

	profiling(cfg.Profile)

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

	cacheDriver, err := redis.NewCacheDriver[any](log, cfg)
	if err != nil {
		return
	}

	trans := translation.NewTranslation(cfg.App.Locale)
	trans.GetLocalizer(cfg.App.Locale)

	queue := messagebroker.NewQueue(log, cfg)
	if err = queue.SetupRabbitMQ(cfg.RabbitMQ.URL, log); err != nil {
		log.Fatal(logger.Queue, logger.Startup, fmt.Sprintf("Failed to setup queue, error: %v", err), nil)
	}

	healthHandler := handler.NewHealthHandler(trans)

	router, err := routes.NewRouter(log, cfg, trans, *healthHandler, cacheDriver)

	userClient := client.NewUserClient(log, cfg.UserManagement)
	defer func() {
		if err = userClient.Close(); err != nil {
			log.Error(logger.Internal, logger.Startup, fmt.Sprintf("Failed to close client connection: %v", err), nil)
		}
	}()
	tokenService := authservice.New(log, cfg.Jwt, cacheDriver)
	otpService := otpservice.New(log, cfg.OTP, cacheDriver)
	oauthService := oauth.New(log, cfg.Oauth)

	permissionRepo := repository.NewPermissionRepository(log, postgresDB)
	roleRepo := repository.NewRoleRepository(log, postgresDB)
	aclRepo := repository.NewACLRepository(log, postgresDB)
	aclService := aclservice.New(log, permissionRepo, roleRepo, aclRepo, userClient)

	authHandler := handler.NewAuthHandler(cfg.OTP, userClient, tokenService, otpService, queue, oauthService, aclService)

	router = router.NewAuthRouter(*authHandler)
	if err != nil {
		log.Error(logger.Internal, logger.Startup, fmt.Sprintf("There is an error when run http: %v", err), nil)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", cfg.Auth.URL, cfg.Auth.Port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router.Handler(),
	}
	log.Info(logger.Internal, logger.Startup, "Starting the HTTP server", map[logger.ExtraKey]interface{}{
		logger.ListeningAddress: server.Addr,
	})

	// Start server
	router.Serve(server)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Server ...", nil)

	timeout := cfg.App.GracefullyShutdown * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal(logger.Internal, logger.Shutdown, fmt.Sprintf("Shutdown Server: %v", err), nil)
	}

	select {
	case <-ctx.Done():
		log.Info(logger.Internal, logger.Shutdown, "timeout of 5 seconds.", nil)
	}
	log.Info(logger.Internal, logger.Shutdown, "Server exiting", nil)
}

func profiling(cfg config.Profile) {
	if cfg.Debug {
		go func() {
			addr := fmt.Sprintf(":%s", cfg.Port)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				return
			}
		}()
	}
}
