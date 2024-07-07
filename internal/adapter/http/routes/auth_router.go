package routes

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/middlewares"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
)

// NewAuthRouter creates a new HTTP router
func (r *Router) NewAuthRouter(
	authHandler handler.AuthHandler,
	roleHandler handler.RoleHandler,
	permissionHandler handler.PermissionHandler,
	authCache port.AuthCache,
) *Router {
	r.Engine.POST("authorize", middlewares.Authentication(r.conf.Jwt, r.trans, authCache), authHandler.Authorize)

	v1 := r.Engine.Group(":language/v1", middlewares.LocaleMiddleware(r.trans))
	{
		auth := v1.Group("auth")
		{
			auth.POST("register", authHandler.Register)
			auth.POST("email-otp/resend", authHandler.EmailOTPResend)
			auth.POST("email-otp/verify", authHandler.EmailOTPVerify)
			auth.POST("login", authHandler.Login)
			auth.POST("google", authHandler.Google)
			auth.POST("forget-password", authHandler.ForgetPassword)
			auth.PATCH("reset-password", authHandler.ResetPassword)
			auth.POST("logout", authHandler.Logout)
		}

		role := v1.Group("roles")
		{
			role.POST("", roleHandler.Create)
			role.GET(":roleID", roleHandler.Get)
			role.GET("", roleHandler.List)
			role.PUT(":roleID", roleHandler.Update)
			role.DELETE(":roleID", roleHandler.Delete)

			role.GET(":roleID/permissions", roleHandler.GetPermissions)
			role.PUT(":roleID/permissions", roleHandler.SyncPermissions)
		}

		v1.GET("permissions", permissionHandler.List)
	}

	return &Router{
		Engine: r.Engine,
		log:    r.log,
		conf:   r.conf,
		trans:  r.trans,
	}
}
