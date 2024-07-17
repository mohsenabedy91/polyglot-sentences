package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/minio"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	trans       translation.Translator
	userService port.UserService
	uowFactory  func() port.UserUnitOfWork
	minioClient *minio.Client
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(
	trans translation.Translator,
	userService port.UserService,
	uowFactory func() port.UserUnitOfWork,
	minioClient *minio.Client,
) *UserHandler {
	return &UserHandler{
		trans:       trans,
		userService: userService,
		uowFactory:  uowFactory,
		minioClient: minioClient,
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
		presenter.NewResponse(ctx, r.trans).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	user, err := r.userService.GetByID(uowFactory, header.UserID)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).ErrorMsg(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, r.trans).Payload(
		presenter.ToUserResource(user),
	).Echo()
}

// Create godoc
// @x-kong {"service": "user-management"}
// @Security AuthBearer[CREATE_USER]
// @Summary Create user
// @Description Create user
// @Tags User
// @Accept json
// @Produce plain
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request formData requests.CreateUserRequest true "Create user request"
// @Param avatar formData file true "Avatar image of the user"
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
		presenter.NewResponse(ctx, r.trans).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var req requests.CreateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		presenter.NewResponse(ctx, r.trans).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := r.userService.IsEmailUnique(uowFactory, req.Email); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	url, uploadFileErr := r.handleFileUpload(ctx, req.Avatar)
	if uploadFileErr != nil {
		presenter.NewResponse(ctx, r.trans).Error(uploadFileErr).Echo(http.StatusInternalServerError)
		return
	}

	user := req.ToUserDomain()
	user.Modifier.CreatedBy = &header.UserID
	user.Avatar = &url

	if _, err := r.userService.Create(uowFactory, user); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, r.trans).Error(err).Echo()
		return
	}

	// TODO send verify Email

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, r.trans).Message(constant.UserSuccessCreate).Echo(http.StatusCreated)
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
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	users, err := r.userService.List(uowFactory)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, r.trans).Error(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, r.trans).Payload(
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
		presenter.NewResponse(ctx, r.trans).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	user, err := r.userService.GetByUUID(uowFactory, userReq.UUIDStr)
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err = uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, r.trans).Payload(
		presenter.ToUserResource(user),
	).Echo()
}

func (r UserHandler) handleFileUpload(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	filePath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(filePath)

	fileTypeSegment := strings.Split(file.Filename, ".")
	uniqueFileName := fmt.Sprintf("%s.%s", uuid.New().String(), fileTypeSegment[len(fileTypeSegment)-1])

	url, uploadFileErr := r.minioClient.UploadFile(ctx.Request.Context(), uniqueFileName, filePath, file.Header.Get("Content-Type"))
	if uploadFileErr != nil {
		return "", uploadFileErr
	}

	return url, nil
}
