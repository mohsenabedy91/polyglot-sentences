package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/claim"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/oauth"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"net/http"
	"sync"
	"time"
)

// AuthHandler represents the HTTP handler for auth-related requests
type AuthHandler struct {
	otpConfig       config.OTP
	userClient      port.UserClient
	tokenService    port.AuthService
	otpCacheService port.OTPCacheService
	queue           *messagebroker.Queue
	oauthService    oauth.GoogleService
	aclService      port.ACLService
	uowFactory      func() repository.UnitOfWork
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(
	otpConfig config.OTP,
	userClient port.UserClient,
	tokenService port.AuthService,
	otpCacheService port.OTPCacheService,
	queue *messagebroker.Queue,
	oauthService oauth.GoogleService,
	aclService port.ACLService,
	uowFactory func() repository.UnitOfWork,
) *AuthHandler {
	return &AuthHandler{
		otpConfig:       otpConfig,
		userClient:      userClient,
		tokenService:    tokenService,
		otpCacheService: otpCacheService,
		queue:           queue,
		oauthService:    oauthService,
		aclService:      aclService,
		uowFactory:      uowFactory,
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
// @ID post_language_v1_auth_register
// @Router /{language}/v1/auth/register [post]
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

	if err := r.otpCacheService.Set(ctx.Request.Context(), req.Email, otp); err != nil {
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

	user, err := r.userClient.Create(ctx.Request.Context(), req.ToDomain())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	uowFactory := r.uowFactory()
	if err = uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err = r.aclService.AssignUserRoleToUser(uowFactory, user.ID); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	message := authservice.SendEmailOTPDto{
		To:       user.Email,
		Name:     user.GetFullName(),
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
// @ID post_language_v1_auth_email_otp_resend
// @Router /{language}/v1/auth/email-otp/resend [post]
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

	if err = r.otpCacheService.Set(ctx.Request.Context(), req.Email, otp); err != nil {
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
// @ID post_language_v1_auth_email_otp_verify
// @Router /{language}/v1/auth/email-otp/verify [post]
func (r AuthHandler) EmailOTPVerify(ctx *gin.Context) {
	var req requests.AuthEmailOTPVerify
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	if err := r.otpCacheService.Validate(ctx.Request.Context(), req.Email, req.Token); err != nil {
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
		_ = r.otpCacheService.Used(ctxWithTimeout, req.Email)

		if !user.WelcomeMessageSent {
			message := authservice.SendWelcomeDto{
				UserID:   user.ID,
				To:       user.Email,
				Name:     user.GetFullName(),
				Language: ctx.Param("language"),
			}
			authservice.SendWelcomeEvent(r.queue, r.userClient).Publish(message)
		}

		if err = r.userClient.UpdateLastLoginTime(ctxWithTimeout, user.ID); err != nil {
			return
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
// @ID post_language_v1_auth_login
// @Router /{language}/v1/auth/login [post]
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

	if ok := helper.CheckPasswordHash(req.Password, *user.Password); !ok {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.CredentialInvalid),
		).Echo()
		return
	}

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = r.userClient.UpdateLastLoginTime(ctxWithTimeout, user.ID); err != nil {
			return
		}
	}()

	token, err := r.tokenService.GenerateToken(user.UUID.String())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToTokenResource(token)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}

// Google godoc
// @Summary Auth Google
// @Description Register or Login Via Google
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.GoogleAuth true "Google request"
// @Success 200 {object} presenter.Response{data=presenter.Token} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_language_v1_auth_google
// @Router /{language}/v1/auth/google [post]
func (r AuthHandler) Google(ctx *gin.Context) {
	var req requests.GoogleAuth
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var (
		googleUserInfo *oauth.GoogleUserInfo
		googleErr      error
		wg             sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		googleUserInfo, googleErr = r.oauthService.UserGoogleInfo(ctx.Request.Context(), req.AccessToken)
	}()

	user, err := r.userClient.GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	wg.Wait()
	if googleErr != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.NewServerError(),
		).Echo()
		return
	}

	if user != nil && user.GoogleID != nil && *user.GoogleID != *googleUserInfo.Id {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.Unauthorized),
		).Echo()
		return
	}

	if user == nil {
		user, err = r.userClient.Create(ctx.Request.Context(), domain.User{
			FirstName: googleUserInfo.FirstName,
			LastName:  googleUserInfo.LastName,
			Email:     googleUserInfo.Email,
			Avatar:    googleUserInfo.AvatarURL,
			GoogleID:  googleUserInfo.Id,
			Status:    domain.UserStatusActive,
		})
		if err != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
			return
		}

		uowFactory := r.uowFactory()
		if err = uowFactory.BeginTx(ctx); err != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
			return
		}

		if err = r.aclService.AssignUserRoleToUser(uowFactory, user.ID); err != nil {
			if rErr := uowFactory.Rollback(); rErr != nil {
				presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
				return
			}
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
			return
		}

		if err = uowFactory.Commit(); err != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
			return
		}

	} else if user.GoogleID == nil {
		if err = r.userClient.UpdateGoogleID(ctx.Request.Context(), user.ID, *googleUserInfo.Id); err != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
			return
		}
	}

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if !user.WelcomeMessageSent {
			message := authservice.SendWelcomeDto{
				UserID:   user.ID,
				To:       user.Email,
				Name:     user.GetFullName(),
				Language: ctx.Param("language"),
			}
			authservice.SendWelcomeEvent(r.queue, r.userClient).Publish(message)
		}
		if err = r.userClient.UpdateLastLoginTime(ctxWithTimeout, user.ID); err != nil {
			return
		}
	}()

	token, err := r.tokenService.GenerateToken(user.UUID.String())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToTokenResource(token)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}

// ForgetPassword godoc
// @Summary Auth ForgetPassword
// @Description Forget Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.ForgetPassword true "ForgetPassword request"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_language_v1_auth_forget_password
// @Router /{language}/v1/auth/forget-password [post]
func (r AuthHandler) ForgetPassword(ctx *gin.Context) {
	var req requests.ForgetPassword
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var (
		otp       string
		otpSetErr error
		wg        sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		otp = helper.GenerateOTP(r.otpConfig.Digits)
		otpSetErr = r.otpCacheService.SetForgetPassword(ctx.Request.Context(), req.Email, otp)
	}()

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

	wg.Wait()
	if otpSetErr != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(otpSetErr).Echo()
		return
	}

	go func() {
		// TODO add rate limit
		message := authservice.SendResetPasswordLinkDto{
			To:       user.Email,
			Name:     user.GetFullName(),
			OTP:      otp,
			Language: ctx.Param("language"),
		}
		authservice.SendResetPasswordLinkEvent(r.queue).Publish(message)
	}()

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessForgetPassword).Echo(http.StatusOK)
}

// ResetPassword godoc
// @Summary Auth ResetPassword
// @Description Reset Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.ResetPassword true "Reset Password request"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID patch_language_v1_auth_reset_password
// @Router /{language}/v1/auth/reset-password [patch]
func (r AuthHandler) ResetPassword(ctx *gin.Context) {
	var req requests.ResetPassword
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var (
		hashedErr    error
		hashPassword string
		wg           sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		hashPassword, hashedErr = helper.HashPassword(req.Password)
	}()

	if err := r.otpCacheService.ValidateForgetPassword(ctx.Request.Context(), req.Email, req.Token); err != nil {
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

	wg.Wait()
	if hashedErr != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(hashedErr).Echo()
		return
	}

	if err = r.userClient.UpdatePassword(ctx.Request.Context(), user.ID, hashPassword); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = r.otpCacheService.UsedForgetPassword(ctxWithTimeout, req.Email)
	}()

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessResetPassword).Echo(http.StatusOK)
}

// Logout godoc
// @Security AuthBearer
// @Summary Logout
// @Description Logout user based on Authorization value
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_language_v1_auth_logout
// @Router /{language}/v1/auth/logout [post]
func (r AuthHandler) Logout(ctx *gin.Context) {
	var header requests.Header
	if err := ctx.ShouldBindHeader(&header); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	if err := r.tokenService.LogoutToken(ctx.Request.Context(), header.JTI, header.EXP); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessLogout).Echo(http.StatusOK)
}

func (r AuthHandler) Authorize(ctx *gin.Context) {
	var req requests.AuthorizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	userUUID := claim.GetUserUUIDFromGinContext(ctx)

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	isAllowed, userID, err := r.aclService.CheckAccess(ctx.Request.Context(), uowFactory, userUUID, req.RequiredPermissions...)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if !isAllowed {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.New(serviceerror.PermissionDenied),
		).Echo()
		return
	}

	data := presenter.Authorize{
		Authorized: isAllowed,
		JTI:        claim.GetJTIFromGinContext(ctx),
		EXP:        claim.GetExpFromGinContext(ctx),
		ID:         userID,
	}

	presenter.NewResponse(ctx, nil).Payload(data).Echo(http.StatusOK)
}
