//go:build !test

package main

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/cmd/setup"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/client"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/routes"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
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
	configProvider := &config.Config{}
	conf := configProvider.GetConfig()
	log := logger.NewLogger(conf.Auth.Name, conf.Log)

	profiling(conf.Profile)

	ctx := context.Background()
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		}
	}()

	postgresDB, err := setup.InitializeDatabase(ctx, log, conf)
	if err != nil {
		return
	}
	uowFactory := func() port.AuthUnitOfWork {
		return repository.NewUnitOfWork(log, postgresDB)
	}

	cache, err := redis.New(log, conf)
	if err != nil {
		return
	}
	defer func() {
		if cacheCloseErr := cache.Close(); cacheCloseErr != nil {
			log.Fatal(logger.Cache, logger.Startup, cacheCloseErr.Error(), nil)
		}
	}()

	queue, err := setup.InitializeQueue(log, conf)
	if err != nil {
		return
	}
	defer queue.Driver.Close()

	trans := translation.NewTranslation(conf.App)
	trans.GetLocalizer(conf.App.Locale)

	userClient := client.NewUserClient(log, conf.UserManagement)
	defer userClient.Close()

	authCache := authrepository.NewAuthCache(log, conf.Redis, cache)
	tokenService := authservice.New(log, conf.Jwt, authCache, nil, nil)

	otpCache := authrepository.NewOTPCache(log, conf.Redis, cache)
	otpCacheService := otpservice.NewOTPCache(conf.OTP, otpCache)

	clientProvider := &oauth.ClientProvider{
		Conf: conf.Oauth,
	}
	oauthService := oauth.New(log, conf.Oauth, clientProvider)

	permissionService := permissionservice.New()

	roleCache := authrepository.NewRoleCache(log, conf.Redis, cache)
	roleCacheService := roleservice.NewRoleCacheService(roleCache)
	roleService := roleservice.New(roleCacheService)

	aclService := aclservice.New(userClient)

	healthHandler := handler.NewHealthHandler(trans)
	authHandler := handler.NewAuthHandler(conf, trans, userClient, tokenService, otpCacheService, queue, oauthService, aclService, uowFactory)
	roleHandler := handler.NewRoleHandler(trans, roleService, uowFactory)
	permissionHandler := handler.NewPermissionHandler(trans, permissionService, uowFactory)

	// Init router
	router, err := routes.NewRouter(log, conf, trans, *healthHandler)
	if err != nil {
		return
	}

	router = router.NewAuthRouter(*authHandler, *roleHandler, *permissionHandler, authCache)

	listenAddr := fmt.Sprintf("%s:%s", conf.Auth.URL, conf.Auth.Port)
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

	timeout := conf.App.GracefullyShutdown * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatal(logger.Internal, logger.Shutdown, fmt.Sprintf("Shutdown Server: %v", err), nil)
	}

	<-ctx.Done()
	log.Info(logger.Internal, logger.Shutdown, "Server exiting", nil)
}

func profiling(conf config.Profile) {
	if conf.Debug {
		go func() {
			addr := fmt.Sprintf(":%s", conf.Port)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				return
			}
		}()
	}
}
