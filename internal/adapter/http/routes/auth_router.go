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
			auth.POST("email-otp/resend", authHandler.EmailOTPResend)
			auth.POST("email-otp/verify", authHandler.EmailOTPVerify)
			auth.POST("login", authHandler.Login)
			auth.GET("profile", middlewares.Authentication(r.cfg.Jwt), authHandler.Profile)
			auth.POST("google", authHandler.Google)
			auth.POST("forget-password", authHandler.ForgetPassword)
			auth.PATCH("reset-password", authHandler.ResetPassword)
		}
	}

	return &Router{
		r.Engine, r.log, r.cfg, r.trans,
	}
}
