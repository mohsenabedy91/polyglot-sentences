package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"net/http"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	userService port.UserService
	uowFactory  func() repository.UnitOfWork
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService port.UserService, uowFactory func() repository.UnitOfWork) *UserHandler {
	return &UserHandler{
		userService: userService,
		uowFactory:  uowFactory,
	}
}

// Profile godoc
// @x-kong {"service": "user-management"}
// @Security AuthBearer
// @Summary Profile
// @Description Get user Profile based on Authorization
// @Tags User
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{data=presenter.User} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_language_v1_users_profile
// @Router /{language}/v1/users/profile [get]
func (r UserHandler) Profile(ctx *gin.Context) {
	var header requests.Header
	if err := ctx.ShouldBindHeader(&header); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	user, err := r.userService.GetByID(uowFactory, header.UserID)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).ErrorMsg(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToUserResource(user)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}

// Create godoc
// @x-kong {"service": "user-management"}
// @Security AuthBearer[CREATE_USER]
// @Summary Create user
// @Description Create user
// @Tags User
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.CreateUserRequest true "Create user request"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_language_v1_users
// @Router /{language}/v1/users [post]
func (r UserHandler) Create(ctx *gin.Context) {
	var header requests.Header
	if err := ctx.ShouldBindHeader(&header); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var req requests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := r.userService.IsEmailUnique(uowFactory, req.Email); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	user := req.ToDomain()
	user.CreatedBy = header.UserID

	if _, err := r.userService.Create(uowFactory, user); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil).Error(err).Echo()
		return
	}

	// TODO send verify Email

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.UserSuccessCreate).Echo(http.StatusCreated)
}

// List godoc
// @x-kong {"service": "user-management"}
// @Security AuthBearer[READ_USER]
// @Summary List of user
// @Description Get list of user
// @Tags User
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{data=[]presenter.User} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_language_v1_users
// @Router /{language}/v1/users [get]
func (r UserHandler) List(ctx *gin.Context) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	users, err := r.userService.List(uowFactory)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil).Error(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Payload(
		presenter.ToUserCollection(users),
	).Echo()
}

// Get godoc
// @x-kong {"service": "user-management"}
// @Security AuthBearer[READ_USER]
// @Summary Get User
// @Description Get User By UUID
// @Tags User
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param userID path string true "user id should be uuid"
// @Success 200 {object} presenter.Response{data=presenter.User} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_language_v1_users_userID
// @Router /{language}/v1/users/{userID} [get]
func (r UserHandler) Get(ctx *gin.Context) {
	var userReq requests.UserUUIDUri
	if err := ctx.ShouldBindUri(&userReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	user, err := r.userService.GetByUUID(uowFactory, userReq.UUIDStr)
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

	presenter.NewResponse(ctx, nil).Payload(
		presenter.ToUserResource(user),
	).Echo()
}
