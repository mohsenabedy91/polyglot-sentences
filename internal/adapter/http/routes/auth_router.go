package routes

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
)

// NewAuthRouter creates a new HTTP router
func (r *Router) NewAuthRouter(authHandler handler.AuthHandler) *Router {
	v1 := r.Engine.Group(":language/v1", middlewares.LocaleMiddleware(r.trans))
	{
		auth := v1.Group("auth")
		{
			auth.POST("register", authHandler.Register)
			auth.POST("login", authHandler.Login)
			auth.GET("profile", middlewares.Authentication(r.config.Jwt), authHandler.Profile)
		}
	}

	return &Router{
		r.Engine, r.log, r.config, r.trans,
	}
}