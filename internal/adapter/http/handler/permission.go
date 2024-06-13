package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"net/http"
)

// PermissionHandler represents the HTTP handler for auth-related requests
type PermissionHandler struct {
	permissionService port.PermissionService
}

// NewPermissionHandler creates a new PermissionHandler instance
func NewPermissionHandler(permissionService port.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// List godoc
// @Security AuthBearer
// @Summary List of Permission
// @Description return a list of permission
// @Tags ACL
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{data=[]presenter.Permission} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_language_v1_permissions
// @Router /{language}/v1/permissions [get]
func (r PermissionHandler) List(ctx *gin.Context) {
	permissions, err := r.permissionService.List(ctx.Request.Context())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	data := presenter.ToPermissionCollection(permissions)

	presenter.NewResponse(ctx, nil).Payload(
		presenter.ToPermissionCollection(permissions),
	).Echo(http.StatusOK)
}
