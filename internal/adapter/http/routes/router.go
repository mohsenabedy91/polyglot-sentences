package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/validations"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Router is a wrapper for HTTP router
type Router struct {
	Engine *gin.Engine
	log    logger.Logger
	conf   config.Config
	trans  translation.Translator
}

// NewRouter creates a new HTTP router
func NewRouter(
	log logger.Logger,
	conf config.Config,
	trans translation.Translator,
	healthHandler handler.HealthHandler,
) (*Router, error) {

	// Disable debug mode in production
	if conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	RegisterPrometheus(log)

	router.Use(middlewares.Prometheus())
	router.Use(gin.Logger(), gin.CustomRecovery(middlewares.ErrorHandler(trans)))
	router.Use(middlewares.DefaultStructuredLogger(log))

	setSwaggerRoutes(router.Group(""), conf.Swagger)

	if err := validations.RegisterValidator(conf); err != nil {
		log.Fatal(logger.Router, logger.Startup, fmt.Sprintf("Failed to setup router, error: %v", err), nil)
		return nil, err
	}

	router.GET("metrics", gin.WrapH(promhttp.Handler()))
	v1 := router.Group(":language/v1", middlewares.LocaleMiddleware(trans))
	{
		v1.GET("health/check", healthHandler.Check)
	}

	return &Router{
		Engine: router,
		log:    log,
		conf:   conf,
		trans:  trans,
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
