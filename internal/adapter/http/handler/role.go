package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"net/http"
)

// RoleHandler represents the HTTP handler for auth-related requests
type RoleHandler struct {
	roleService port.RoleService
	uowFactory  func() repository.UnitOfWork
}

// NewRoleHandler creates a new RoleHandler instance
func NewRoleHandler(roleService port.RoleService, uowFactory func() repository.UnitOfWork) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		uowFactory:  uowFactory,
	}
}

// Create godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[CREATE_ROLE]
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
// @ID post_language_v1_roles
// @Router /{language}/v1/roles [post]
func (r RoleHandler) Create(ctx *gin.Context) {
	var req requests.RoleCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	role := domain.Role{
		Title:       req.Title,
		Description: req.Description,
	}
	if err := r.roleService.Create(uowFactory, role); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.RoleSuccessCreated).Echo(http.StatusCreated)
}

// Get godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[READ_ROLE]
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
// @ID get_language_v1_roles_roleID
// @Router /{language}/v1/roles/{roleID} [get]
func (r RoleHandler) Get(ctx *gin.Context) {
	var userReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&userReq); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	role, err := r.roleService.Get(uowFactory, userReq.UUIDStr)
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
		presenter.ToRoleResource(role),
	).Echo(http.StatusOK)
}

// List godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[READ_ROLE]
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
// @ID get_language_v1_roles
// @Router /{language}/v1/roles [get]
func (r RoleHandler) List(ctx *gin.Context) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	roles, err := r.roleService.List(ctx.Request.Context(), uowFactory)
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
		presenter.ToRoleCollection(roles),
	).Echo(http.StatusOK)
}

// Update godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[UPDATE_ROLE]
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
// @ID put_language_v1_roles_roleID
// @Router /{language}/v1/roles/{roleID} [put]
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

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	role := domain.Role{
		Title:       req.Title,
		Description: req.Description,
	}
	if err := r.roleService.Update(ctx.Request.Context(), uowFactory, role, roleReq.UUIDStr); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.RoleSuccessUpdated).Echo(http.StatusOK)
}

// Delete godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[DELETE_ROLE]
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
// @ID delete_language_v1_roles_roleID
// @Router /{language}/v1/roles/{roleID} [delete]
func (r RoleHandler) Delete(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := r.roleService.Delete(uowFactory, roleReq.UUIDStr); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Echo(http.StatusNoContent)
}

// GetPermissions godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[READ_ROLE_PERMISSIONS]
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
// @ID get_language_v1_roles_roleID_permissions
// @Router /{language}/v1/roles/{roleID}/permissions [get]
func (r RoleHandler) GetPermissions(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	rolePermissions, err := r.roleService.GetPermissions(uowFactory, roleReq.UUIDStr)
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
		presenter.ToRoleResource(rolePermissions),
	).Echo(http.StatusOK)
}

// SyncPermissions godoc
// @x-kong {"service": "auth"}
// @Security AuthBearer[SYNC_PERMISSIONS_WITH_ROLE]
// @Summary Sync Permissions
// @Description Assign/Remove permissions for a Role
// @Tags Role
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Param roleID path string true "role id should be uuid"
// @Param request body requests.SyncPermissions true "Assign Permissions"
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Failure 400 {object} presenter.Error "Failed response"
// @Failure 401 {object} presenter.Error "Unauthorized"
// @Failure 422 {object} presenter.Response{validationErrors=[]presenter.ValidationError} "Validation error"
// @Failure 500 {object} presenter.Error "Internal server error"
// @ID put_language_v1_roles_roleID_permissions
// @Router /{language}/v1/roles/{roleID}/permissions [put]
func (r RoleHandler) SyncPermissions(ctx *gin.Context) {
	var roleReq requests.RoleUUIDUri
	if err := ctx.ShouldBindUri(&roleReq); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	var req requests.SyncPermissions
	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.NewResponse(ctx, nil).Validation(err).Echo(http.StatusUnprocessableEntity)
		return
	}

	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := r.roleService.SyncPermissions(uowFactory, roleReq.UUIDStr, req.Permissions); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(rErr).Echo()
			return
		}
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	if err := uowFactory.Commit(); err != nil {
		presenter.NewResponse(ctx, nil, StatusCodeMapping).Error(err).Echo()
		return
	}

	presenter.NewResponse(ctx, nil).Message(constant.RoleSuccessSyncPermissions).Echo(http.StatusOK)
}
