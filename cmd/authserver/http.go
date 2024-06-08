package main

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/cmd/setup"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/client"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/routes"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/repository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/aclservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/otpservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/permissionservice"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/roleservice"
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

	cacheDriver, err := redis.NewCacheDriver[any](log, cfg)
	if err != nil {
		return
	}
	defer cacheDriver.Close()

	queue, err := setup.InitializeQueue(log, cfg)
	if err != nil {
		return
	}
	defer queue.Driver.Close()

	trans := translation.NewTranslation(cfg.App.Locale)
	trans.GetLocalizer(cfg.App.Locale)

	userClient := client.NewUserClient(log, cfg.UserManagement)
	defer userClient.Close()

	tokenService := authservice.New(log, cfg.Jwt, cacheDriver)

	otpCacheService := otpservice.NewOTPCache(log, cfg.OTP, cacheDriver)

	oauthService := oauth.New(log, cfg.Oauth)

	permissionRepo := repository.NewPermissionRepository(log, postgresDB)
	permissionService := permissionservice.New(log, permissionRepo)

	roleRepo := repository.NewRoleRepository(log, postgresDB)
	roleCache := roleservice.NewRoleCache(log, cacheDriver)
	roleService := roleservice.New(log, roleRepo, roleCache)

	aclRepo := repository.NewACLRepository(log, postgresDB)
	aclService := aclservice.New(log, permissionRepo, roleRepo, aclRepo, userClient)

	healthHandler := handler.NewHealthHandler(trans)
	authHandler := handler.NewAuthHandler(cfg.OTP, userClient, tokenService, otpCacheService, queue, oauthService, aclService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// Init router
	router, err := routes.NewRouter(log, cfg, trans, cacheDriver, aclService, *healthHandler)
	if err != nil {
		return
	}

	router = router.NewAuthRouter(*authHandler, *roleHandler, *permissionHandler)

	listenAddr := fmt.Sprintf("%s:%s", cfg.Auth.URL, cfg.Auth.Port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router.Engine.Handler(),
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
