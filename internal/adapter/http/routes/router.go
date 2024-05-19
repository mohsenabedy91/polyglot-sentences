package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/validations"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
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
			user.GET(":userID", userHandler.Get)
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
