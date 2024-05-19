package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"net/http"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	userService port.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userSvc port.UserService) *UserHandler {
	return &UserHandler{
		userService: userSvc,
	}
}

// Get godoc
// @Security AuthBearer
// @Summary Get User
// @Description Get User By UUID
// @Tags User
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param userID path string true "user id should be uuid"
// @Success 200 {object} presenter.Response{data=presenter.User} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_v1_users_userID
// @Router /v1/users/{userID} [get]
func (r UserHandler) Get(ctx *gin.Context) {
	var req requests.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	user, err := r.userService.GetByUUID(ctx.Request.Context(), uuid.MustParse(req.UserID))
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	result := presenter.ToUserResource(user)

	presenter.NewResponse(ctx, nil).Payload(result).Echo()
}
