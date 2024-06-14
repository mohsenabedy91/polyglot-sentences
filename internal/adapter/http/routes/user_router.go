package routes

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
)

// NewUserRouter creates a new HTTP router
func (r *Router) NewUserRouter(userHandler handler.UserHandler) *Router {
	v1 := r.Engine.Group(":language/v1", middlewares.LocaleMiddleware(r.trans))
	{
		user := v1.Group("users")
		{
			user.GET("profile", userHandler.Profile)
			user.POST("", userHandler.Create)
			user.GET("", userHandler.List)
			user.GET(":userID", userHandler.Get)
		}
	}

	return &Router{
		Engine: r.Engine,
		log:    r.log,
		cfg:    r.cfg,
		trans:  r.trans,
		cache:  r.cache,
	}
}
