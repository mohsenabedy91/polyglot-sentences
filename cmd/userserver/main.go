//go:build !test

package main

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/cmd/setup"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/server"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/routes"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/minio"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"google.golang.org/grpc"
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
	configProvider := &config.Config{}
	conf := configProvider.GetConfig()
	log := logger.NewLogger(conf.UserManagement.Name, conf.Log)

	ctx := context.Background()
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		}
	}()
	postgresDB, err := setup.InitializeDatabase(ctx, log, conf)
	if err != nil {
		log.Fatal(logger.Database, logger.Startup, err.Error(), nil)
		return
	}
	uowFactory := func() repository.UnitOfWork {
		return repository.NewUnitOfWork(log, postgresDB)
	}

	trans := translation.NewTranslation(conf.App)
	trans.GetLocalizer(conf.App.Locale)

	userService := userservice.New(log)

	grpcServer := startGRPCServer(conf, log, userService, uowFactory)
	httpServer := startHTTPServer(ctx, conf, log, userService, uowFactory, trans)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Servers ...", nil)

	grpcServer.GracefulStop()
	shutdownHTTPServer(ctx, httpServer, log, conf)
}

func startGRPCServer(
	conf config.Config,
	log logger.Logger,
	userService *userservice.UserService,
	uowFactory func() repository.UnitOfWork,
) *grpc.Server {
	s := server.NewUserGRPCServer(conf.UserManagement, userService, uowFactory)
	grpcServer, err := s.StartUserGRPCServer()
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return nil
	}

	log.Info(logger.Internal, logger.Startup, "gRPC server started", nil)
	return grpcServer
}

func startHTTPServer(
	ctx context.Context,
	conf config.Config,
	log logger.Logger,
	userService *userservice.UserService,
	uowFactory func() repository.UnitOfWork,
	trans *translation.Translation,
) *http.Server {
	cacheDriver, err := redis.NewCacheDriver[any](log, conf)
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return nil
	}
	defer cacheDriver.Close()

	minioClient, err := minio.NewMinioClient(ctx, log, conf.Minio)
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return nil
	}

	userHandler := handler.NewUserHandler(userService, uowFactory, minioClient)
	healthHandler := handler.NewHealthHandler(trans)

	// Init router
	router, err := routes.NewRouter(log, conf, trans, cacheDriver, *healthHandler)
	if err != nil {
		log.Fatal(logger.Internal, logger.Startup, err.Error(), nil)
		return nil
	}

	router = router.NewUserRouter(*userHandler)

	listenAddr := fmt.Sprintf("%s:%s", conf.UserManagement.URL, conf.UserManagement.HTTPPort)
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: router.Engine.Handler(),
	}
	log.Info(logger.Internal, logger.Startup, "Starting the HTTP server", map[logger.ExtraKey]interface{}{
		logger.ListeningAddress: httpServer.Addr,
	})

	go router.Serve(httpServer)
	return httpServer
}

func shutdownHTTPServer(
	ctx context.Context,
	server *http.Server,
	log logger.Logger,
	conf config.Config,
) {
	timeout := conf.App.GracefullyShutdown * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(logger.Internal, logger.Shutdown, fmt.Sprintf("Shutdown Server: %v", err), nil)
	}

	select {
	case <-ctx.Done():
		log.Info(logger.Internal, logger.Shutdown, "timeout of 5 seconds.", nil)
	}
	log.Info(logger.Internal, logger.Shutdown, "Server exiting", nil)
}
