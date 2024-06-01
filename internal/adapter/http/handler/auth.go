package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/claim"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"net/http"
	"sync"
	"time"
)

// AuthHandler represents the HTTP handler for auth-related requests
type AuthHandler struct {
	otpConfig    config.OTP
	userClient   port.UserClient
	tokenService port.AuthService
	otpService   port.OtpService
	queue        *messagebroker.Queue
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(
	otpConfig config.OTP,
	userClient port.UserClient,
	tokenService port.AuthService,
	otpService port.OtpService,
	queue *messagebroker.Queue,
) *AuthHandler {
	return &AuthHandler{
		otpConfig:    otpConfig,
		userClient:   userClient,
		tokenService: tokenService,
		otpService:   otpService,
		queue:        queue,
	}
}

// Register godoc
// @Summary Auth Register
// @Description Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.AuthRegister true "Register request"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_v1_auth_register
// @Router /v1/auth/register [post]
func (r AuthHandler) Register(ctx *gin.Context) {
	var req requests.AuthRegister
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var (
		hashedPass string
		hashErr    error
		wg         sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		hashedPass, hashErr = helper.HashPassword(req.Password)
	}()

	if err := r.userClient.IsEmailUnique(ctx.Request.Context(), req.Email); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	otp := helper.GenerateOTP(r.otpConfig.Digits)

	if err := r.otpService.Set(ctx.Request.Context(), req.Email, otp); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}
	wg.Wait()

	if hashErr != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.NewServerError(),
		).Echo()
		return
	}

	req.Password = hashedPass

	if err := r.userClient.Create(ctx.Request.Context(), req.ToDomain()); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	message := authservice.SendEmailOTPDto{
		To:       req.Email,
		Name:     req.FirstName + " " + req.LastName,
		OTP:      otp,
		Language: ctx.Param("language"),
	}
	authservice.SendEmailOTPEvent(r.queue).Publish(message)

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessRegisteredUser).Echo(http.StatusCreated)
}

// EmailOTPResend godoc
// @Summary EmailOTPResend
// @Description Resend OTP to User via Email
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.AuthEmailOTPResend true "EmailOTPResend request"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_v1_auth_email_otp_resend
// @Router /v1/auth/email-otp/resend [post]
func (r AuthHandler) EmailOTPResend(ctx *gin.Context) {
	var req requests.AuthEmailOTPResend
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	user, err := r.userClient.GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}
	if user == nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.RecordNotFound),
		).Echo()
		return
	}

	otp := helper.GenerateOTP(r.otpConfig.Digits)

	if err = r.otpService.Set(ctx.Request.Context(), req.Email, otp); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	// TODO add rate limit
	message := authservice.SendEmailOTPDto{
		To:       user.Email,
		Name:     user.GetFullName(),
		OTP:      otp,
		Language: ctx.Param("language"),
	}
	authservice.SendEmailOTPEvent(r.queue).Publish(message)

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessEmailOTPSent).Echo(http.StatusOK)
}

// EmailOTPVerify godoc
// @Summary EmailOTPVerify
// @Description Verify User via Email OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.AuthEmailOTPVerify true "EmailOTPVerify request"
// @Success 200 {object} presenter.Response{data=presenter.Token} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_v1_auth_email_otp_verify
// @Router /v1/auth/email-otp/verify [post]
func (r AuthHandler) EmailOTPVerify(ctx *gin.Context) {
	var req requests.AuthEmailOTPVerify
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	if err := r.otpService.Validate(ctx.Request.Context(), req.Email, req.Token); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := r.userClient.VerifiedEmail(ctx.Request.Context(), req.Email); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	user, err := r.userClient.GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}
	if user == nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.RecordNotFound),
		).Echo()
		return
	}

	token, err := r.tokenService.GenerateToken(user.UUID.String())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		_ = r.otpService.Used(ctxWithTimeout, req.Email)

		if !user.WelcomeMessageSent {
			message := authservice.SendWelcomeDto{
				To:       user.Email,
				Name:     user.GetFullName(),
				Language: ctx.Param("language"),
			}
			authservice.SendWelcomeEvent(r.queue).Publish(message)

			if err = r.userClient.MarkWelcomeMessageSent(ctxWithTimeout, user.ID); err != nil {
				return
			}
		}
	}()

	result := presenter.ToTokenResource(token)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}

// Login godoc
// @Summary Login
// @Description User based on email and password can log in app
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.AuthLogin true "Login request"
// @Success 200 {object} presenter.Response{data=presenter.Token} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_v1_auth_login
// @Router /v1/auth/login [post]
func (r AuthHandler) Login(ctx *gin.Context) {
	var req requests.AuthLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	user, err := r.userClient.GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}
	if user == nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.RecordNotFound),
		).Echo()
		return
	}

	if ok := helper.CheckPasswordHash(req.Password, user.Password); !ok {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.CredentialInvalid),
		).Echo()
		return
	}

	token, err := r.tokenService.GenerateToken(user.UUID.String())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToTokenResource(token)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}

// Profile godoc
// @Security AuthBearer
// @Summary Profile
// @Description Get user Profile based on Authorization
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{data=presenter.User} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_v1_auth_profile
// @Router /v1/auth/profile [get]
func (r AuthHandler) Profile(ctx *gin.Context) {
	userUUID := claim.GetUserUUIDFromGinContext(ctx)

	user, err := r.userClient.GetByUUID(ctx.Request.Context(), userUUID.String())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).ErrorMsg(err).Echo()
		return
	}

	result := presenter.ToUserResource(user)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}
