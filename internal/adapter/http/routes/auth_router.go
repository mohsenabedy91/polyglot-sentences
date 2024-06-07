package routes

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

// NewAuthRouter creates a new HTTP router
func (r *Router) NewAuthRouter(
	authHandler handler.AuthHandler,
	roleHandler handler.RoleHandler,
) *Router {
	v1 := r.Engine.Group(":language/v1", middlewares.LocaleMiddleware(r.trans))
	{
		auth := v1.Group("auth")
		{
			auth.POST("register", authHandler.Register)
			auth.POST("email-otp/resend", authHandler.EmailOTPResend)
			auth.POST("email-otp/verify", authHandler.EmailOTPVerify)
			auth.POST("login", authHandler.Login)
			auth.GET("profile", middlewares.Authentication(r.cfg.Jwt, r.cache), authHandler.Profile)
			auth.POST("google", authHandler.Google)
			auth.POST("forget-password", authHandler.ForgetPassword)
			auth.PATCH("reset-password", authHandler.ResetPassword)
			auth.POST("logout", middlewares.Authentication(r.cfg.Jwt, r.cache), authHandler.Logout)
		}

		role := v1.Group("roles", middlewares.Authentication(r.cfg.Jwt, r.cache))
		{
			role.POST("", middlewares.ACL(r.aclService, domain.PermissionKeyCreateRole), roleHandler.Create)
			role.GET(":roleID", middlewares.ACL(r.aclService, domain.PermissionKeyReadRole), roleHandler.Get)
			role.GET("", middlewares.ACL(r.aclService, domain.PermissionKeyReadRole), roleHandler.List)
			role.PUT(":roleID", middlewares.ACL(r.aclService, domain.PermissionKeyUpdateRole), roleHandler.Update)
			role.DELETE(":roleID", middlewares.ACL(r.aclService, domain.PermissionKeyDeleteRole), roleHandler.Delete)
		}
	}

	return &Router{
		Engine:     r.Engine,
		log:        r.log,
		cfg:        r.cfg,
		aclService: r.aclService,
		trans:      r.trans,
		cache:      r.cache,
	}
}
