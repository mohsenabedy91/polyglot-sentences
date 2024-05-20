package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/claim"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/password"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"net/http"
)

// AuthHandler represents the HTTP handler for auth-related requests
type AuthHandler struct {
	userService  port.UserService
	tokenService port.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(userService port.UserService, tokenService port.AuthService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
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

	hashedPass := make(chan string)
	go func() {
		err := password.HashPassword(req.Password, hashedPass)
		if err != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
				serviceerror.NewServiceError(serviceerror.ServerError),
			).Echo()
			return
		}
	}()

	if err := r.userService.IsEmailUnique(ctx.Request.Context(), req.Email); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}
	req.Password = <-hashedPass

	if err := r.userService.RegisterUser(ctx.Request.Context(), req.ToDomain()); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.AuthSuccessRegisteredUser).Echo(http.StatusCreated)
}

// Login godoc
// @Summary Login
// @Description User based on email and password can log in app
// @Tags Auth
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.AuthLogin true "Login request"
// @Success 200 {object} presenter.Response{data=presenter.User} "Successful response"
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

	user, err := r.userService.GetByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if ok := password.CheckPasswordHash(req.Password, *user.Password); !ok {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(
			serviceerror.NewServiceError(serviceerror.CredentialInvalid),
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

	user, err := r.userService.GetByUUID(ctx.Request.Context(), userUUID)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToUserResource(user)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}
