package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"net/http"
)

// PermissionHandler represents the HTTP handler for auth-related requests
type PermissionHandler struct {
	trans             translation.Translator
	permissionService port.PermissionService
	uowFactory        func() port.AuthUnitOfWork
}

// NewPermissionHandler creates a new PermissionHandler instance
func NewPermissionHandler(
	trans translation.Translator,
	permissionService port.PermissionService,
	uowFactory func() port.AuthUnitOfWork,
) *PermissionHandler {
	return &PermissionHandler{
		trans:             trans,
		permissionService: permissionService,
		uowFactory:        uowFactory,
	}
}

// List godoc
// @x-kong {"service": "auth-service"}
// @Security AuthBearer[READ_PERMISSION]
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
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, r.trans, StatusCodeMapping).Error(err).Echo()
		return
	}

	permissions, err := r.permissionService.List(uowFactory)
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
		presenter.ToPermissionCollection(permissions),
	).Echo(http.StatusOK)
}
