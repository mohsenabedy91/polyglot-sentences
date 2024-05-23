package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/validations"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"net/http"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
	log logger.Logger
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.Config,
	log logger.Logger,
	trans *translation.Translation,
	accessControlService port.AccessControlService,
	healthHandler handler.HealthHandler,
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
) (*Router, error) {

	// Disable debug mode in production
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	RegisterPrometheus(log)

	router.Use(gin.Logger(), gin.CustomRecovery(middlewares.ErrorHandler))
	router.Use(middlewares.DefaultStructuredLogger(log))

	setSwaggerRoutes(router.Group(""), config.Swagger)

	if err := validations.RegisterValidator(); err != nil {
		return nil, err
	}

	v1 := router.Group(":language/v1", middlewares.LocaleMiddleware(trans))
	{
		v1.GET("health/check", healthHandler.Check)
		auth := v1.Group("auth")
		{
			auth.POST("register", authHandler.Register)
			auth.POST("login", authHandler.Login)
			auth.GET("profile", middlewares.Authentication(config.Jwt), authHandler.Profile)
		}
		user := v1.Group("users", middlewares.Authentication(config.Jwt), middlewares.ACL(
			accessControlService,
			domain.PermissionKeyCreateUser,
			domain.PermissionKeyReadUser,
		))
		{
			user.POST("", userHandler.Create)
			user.GET("", userHandler.List)
			user.GET(":userID", userHandler.Get)
		}
	}

	return &Router{
		router, log,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(server *http.Server) {
	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			r.log.Error(logger.Internal, logger.Startup, fmt.Sprintf("Error starting the HTTP server: %v", err), nil)
		}
	}()
}

func RegisterPrometheus(log logger.Logger) {
	err := prometheus.Register(metrics.DbCall)
	if err != nil {
		log.Error(logger.Prometheus, logger.Startup, err.Error(), nil)
	}

	err = prometheus.Register(metrics.HttpDuration)
	if err != nil {
		log.Error(logger.Prometheus, logger.Startup, err.Error(), nil)
	}
}
