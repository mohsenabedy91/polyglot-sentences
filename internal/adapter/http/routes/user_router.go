package routes

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
)

// NewUserRouter creates a new HTTP router
func (r *Router) NewUserRouter(
	aclService port.ACLService,
	userHandler handler.UserHandler,
) *Router {
	v1 := r.Engine.Group(":language/v1", middlewares.LocaleMiddleware(r.trans))
	{
		user := v1.Group(
			"users",
			middlewares.Authentication(r.cfg.Jwt, r.cache),
			middlewares.ACL(
				aclService,
				domain.PermissionKeyCreateUser,
				domain.PermissionKeyReadUser,
			),
		)
		{
			user.POST("", userHandler.Create)
			user.GET("", userHandler.List)
			user.GET(":userID", userHandler.Get)
		}
	}

	return &Router{
		r.Engine, r.log, r.cfg, r.trans, r.cache,
	}
}
