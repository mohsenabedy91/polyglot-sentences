package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"net/http"
)

// RoleHandler represents the HTTP handler for auth-related requests
type RoleHandler struct {
	roleService port.RoleService
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(roleService port.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Create godoc
// @Security AuthBearer
// @Summary Create a Role
// @Description Create a Role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param request body requests.RoleCreate true "Role Create"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID post_v1_roles
// @Router /v1/roles [post]
func (r RoleHandler) Create(ctx *gin.Context) {
	var req requests.RoleCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	err := r.roleService.Create(ctx.Request.Context(), domain.Role{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.RoleSuccessCreated).Echo(http.StatusCreated)
}

// Get godoc
// @Security AuthBearer
// @Summary Get a Role
// @Description return a role by role uuid
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param roleID path string true "role id should be uuid"
// @Success 200 {object} presenter.Response{data=presenter.Role} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 404 {object} presenter.Error "Not found"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_v1_roles_roleID
// @Router /v1/roles/{roleID} [get]
func (r RoleHandler) Get(ctx *gin.Context) {
	var userReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&userReq); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	role, err := r.roleService.Get(ctx.Request.Context(), userReq.UUIDStr)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	data := presenter.ToRoleResource(role)

	presenter.NewResponse(ctx, nil).Payload(data).Echo(http.StatusOK)
}

// List godoc
// @Security AuthBearer
// @Summary List of Role
// @Description return a list of role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{data=[]presenter.Role} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_v1_roles
// @Router /v1/roles [get]
func (r RoleHandler) List(ctx *gin.Context) {
	roles, err := r.roleService.List(ctx.Request.Context())
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	data := presenter.ToRoleCollection(roles)

	presenter.NewResponse(ctx, nil).Payload(data).Echo(http.StatusOK)
}

// Update godoc
// @Security AuthBearer
// @Summary Update a Role
// @Description Update a Role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param roleID path string true "role id should be uuid"
// @Param request body requests.RoleUpdate true "Create Role"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID put_v1_roles_roleID
// @Router /v1/roles/{roleID} [put]
func (r RoleHandler) Update(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}
	var req requests.RoleUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	err := r.roleService.Update(ctx.Request.Context(), domain.Role{
		Title:       req.Title,
		Description: req.Description,
	}, roleReq.UUIDStr)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.RoleSuccessUpdated).Echo(http.StatusOK)
}

// Delete godoc
// @Security AuthBearer
// @Summary Delete a Role
// @Description Delete a Role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param roleID path string true "role id should be uuid"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 403 {object} presenter.Error "Forbidden"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID delete_v1_roles_roleID
// @Router /v1/roles/{roleID} [delete]
func (r RoleHandler) Delete(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	err := r.roleService.Delete(ctx.Request.Context(), roleReq.UUIDStr)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Echo(http.StatusNoContent)
}

// GetPermissions godoc
// @Security AuthBearer
// @Summary Get Permissions
// @Description get permissions for a Role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param roleID path string true "role id should be uuid"
// @Success 200 {object} presenter.Response{data=presenter.Role} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID get_v1_roles_roleID_permissions
// @Router /v1/roles/{roleID}/permissions [get]
func (r RoleHandler) GetPermissions(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	rolePermissions, err := r.roleService.GetPermissions(ctx.Request.Context(), roleReq.UUIDStr)
	if err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	data := presenter.ToRoleResource(rolePermissions)

	presenter.NewResponse(ctx, nil).Payload(data).Echo(http.StatusOK)
}
