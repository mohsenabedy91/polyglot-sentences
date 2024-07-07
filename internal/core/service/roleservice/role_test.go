package roleservice_test

import (
	"context"
	"github.com/bxcodec/faker/v4"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/roleservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

var wg sync.WaitGroup

func TestRoleService_Create(t *testing.T) {
	roleID := uuid.New()
	role := domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title: "Admin",
		Key:   "ADMIN",
	}

	t.Run("Create success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("ExistKey", role.Key).Return(false, nil)
		mockRepo.On("Create", role).Return(nil)

		service := roleservice.New(nil)
		err := service.Create(mockUow, role)

		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Create role exists error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("ExistKey", role.Key).Return(true, nil)

		service := roleservice.New(nil)
		err := service.Create(mockUow, role)

		require.Error(t, err)
		require.Equal(t, serviceerror.RoleExisted, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_Get(t *testing.T) {
	roleID := uuid.New()
	role := &domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title: "Admin",
		Key:   "ADMIN",
	}

	t.Run("Get success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("GetByUUID", roleID).Return(role, nil)

		service := roleservice.New(nil)
		result, err := service.Get(mockUow, roleID.String())

		require.NoError(t, err)
		require.Equal(t, role, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Get error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("GetByUUID", roleID).Return(&domain.Role{}, serviceerror.New(serviceerror.RecordNotFound))

		service := roleservice.New(nil)
		result, err := service.Get(mockUow, roleID.String())

		require.Error(t, err)
		require.Equal(t, serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
		require.Equal(t, &domain.Role{}, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestRoleService_List(t *testing.T) {
	roleID := uuid.New()
	defaultRole := &domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title:     "Admin",
		Key:       "ADMIN",
		IsDefault: true,
	}

	nonDefaultRole := &domain.Role{
		Base: domain.Base{
			UUID: uuid.New(),
		},
		Title:     "User",
		Key:       "USER",
		IsDefault: false,
	}

	ctx := context.TODO()

	t.Run("List success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		roles := []*domain.Role{defaultRole, nonDefaultRole}
		mockRepo.On("List").Return(roles, nil)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)

		wg.Add(1)
		mockRoleCacheService.On("SetBulk", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				defer wg.Done()
				cacheRoles := args.Get(1).(map[string]domain.RoleKeyType)
				require.Equal(t, defaultRole.Key, cacheRoles[defaultRole.UUID.String()])
				_, exists := cacheRoles[nonDefaultRole.UUID.String()]
				require.False(t, exists)
			}).
			Return(nil)

		service := roleservice.New(mockRoleCacheService)
		service.SetRoleCache(mockRoleCacheService)

		result, err := service.List(ctx, mockUow)

		wg.Wait()

		require.NoError(t, err)
		require.Equal(t, roles, result)

		mockRepo.AssertExpectations(t)
		mockRoleCacheService.AssertExpectations(t)
	})

	t.Run("List error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("List").Return([]*domain.Role{}, serviceerror.NewServerError())

		service := roleservice.New(nil)
		result, err := service.List(ctx, mockUow)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
		require.Equal(t, []*domain.Role{}, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Cache set error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		roles := []*domain.Role{defaultRole}
		mockRepo.On("List").Return(roles, nil)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)

		wg.Add(1)
		mockRoleCacheService.On("SetBulk", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				defer wg.Done()
			}).
			Return(serviceerror.NewServerError())

		service := roleservice.New(mockRoleCacheService)
		service.SetRoleCache(mockRoleCacheService)
		result, err := service.List(ctx, mockUow)

		wg.Wait()

		require.NoError(t, err)
		require.Equal(t, roles, result)

		mockRepo.AssertExpectations(t)
		mockRoleCacheService.AssertExpectations(t)
	})
}

func TestRoleService_Update(t *testing.T) {
	roleID := uuid.New()
	role := domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title: "Admin",
		Key:   "ADMIN",
	}

	ctx := context.TODO()

	t.Run("Update success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("ExistKey", role.Key).Return(false, nil)
		mockRepo.On("Update", role, roleID).Return(nil)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)
		var roleKey = &role.Key
		mockRoleCacheService.On("Get", ctx, roleID.String()).Return(roleKey, nil)

		service := roleservice.New(mockRoleCacheService)
		err := service.Update(ctx, mockUow, role, roleID.String())

		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockRoleCacheService.AssertExpectations(t)
	})

	t.Run("Update success cache return nil", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("ExistKey", role.Key).Return(false, nil)
		mockRepo.On("Update", role, roleID).Return(nil)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)
		var roleKey *domain.RoleKeyType
		mockRoleCacheService.On("Get", ctx, roleID.String()).Return(roleKey, nil)

		service := roleservice.New(mockRoleCacheService)
		err := service.Update(ctx, mockUow, role, roleID.String())

		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockRoleCacheService.AssertExpectations(t)
	})

	t.Run("Update failed key exist", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("ExistKey", role.Key).Return(true, nil)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)
		var roleKey *domain.RoleKeyType
		mockRoleCacheService.On("Get", ctx, roleID.String()).Return(roleKey, nil)

		service := roleservice.New(mockRoleCacheService)
		err := service.Update(ctx, mockUow, role, roleID.String())

		require.Error(t, err)
		require.Equal(t, serviceerror.RoleExisted, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRepo.AssertExpectations(t)
		mockRoleCacheService.AssertExpectations(t)
	})

	t.Run("Update failed cache error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRoleCacheService := new(roleservice.MockRoleCacheService)
		var roleKey *domain.RoleKeyType
		mockRoleCacheService.On("Get", ctx, roleID.String()).Return(roleKey, serviceerror.New(serviceerror.RecordNotFound))

		service := roleservice.New(mockRoleCacheService)
		err := service.Update(ctx, mockUow, role, roleID.String())

		require.Error(t, err)
		require.Equal(t, serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRoleCacheService.AssertExpectations(t)
	})
}

func TestRoleService_Delete(t *testing.T) {
	roleID := uuid.New()

	t.Run("Delete success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("Delete", roleID, uint64(1)).Return(nil)

		service := roleservice.New(nil)
		err := service.Delete(mockUow, roleID.String(), uint64(1))

		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("Delete", roleID, uint64(1)).Return(serviceerror.New(serviceerror.IsNotDeletable))

		service := roleservice.New(nil)
		err := service.Delete(mockUow, roleID.String(), uint64(1))

		require.Error(t, err)
		require.Equal(t, serviceerror.IsNotDeletable, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetPermissions(t *testing.T) {
	roleID := uuid.New()
	role := &domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title: "Admin",
		Key:   "ADMIN",
	}

	t.Run("GetPermissions success", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("GetPermissions", roleID).Return(role, nil)

		service := roleservice.New(nil)
		result, err := service.GetPermissions(mockUow, roleID.String())

		require.NoError(t, err)
		require.Equal(t, role, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPermissions error", func(t *testing.T) {
		mockRepo := new(authrepository.MockRoleRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRepo)

		mockRepo.On("GetPermissions", roleID).Return(&domain.Role{}, serviceerror.New(serviceerror.RecordNotFound))

		service := roleservice.New(nil)
		result, err := service.GetPermissions(mockUow, roleID.String())

		require.Error(t, err)
		require.Equal(t, serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
		require.Equal(t, &domain.Role{}, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestService_SyncPermissions(t *testing.T) {
	roleID := uuid.New()
	role := &domain.Role{
		Base: domain.Base{
			UUID: roleID,
		},
		Title: "Admin",
		Key:   "ADMIN",
	}

	var permissions []domain.Permission
	for i := 1; i <= 5; i++ {
		permissions = append(permissions, domain.Permission{
			Base: domain.Base{
				UUID: uuid.New(),
				ID:   uint64(i),
			},
			Title: helper.StringPtr(faker.Word()),
		})
	}

	t.Run("SyncPermissions success", func(t *testing.T) {
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRoleRepo)
		mockUow.On("PermissionRepository").Return(mockPermissionRepo)

		var permissionUUIDStr []string
		var permissionUUIDs []uuid.UUID
		var validPermissionIDs []uint64
		for _, permission := range permissions {
			permissionUUIDStr = append(permissionUUIDStr, permission.UUID.String())
			permissionUUIDs = append(permissionUUIDs, permission.UUID)
			validPermissionIDs = append(validPermissionIDs, permission.ID)
		}

		mockRoleRepo.On("GetByUUID", roleID).Return(role, nil)
		mockPermissionRepo.On("FilterValidPermissions", permissionUUIDs).Return(validPermissionIDs, nil)
		mockRoleRepo.On("SyncPermissions", role.ID, validPermissionIDs).Return(nil)

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, roleID.String(), permissionUUIDStr)

		require.NoError(t, err)

		mockRoleRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("SyncPermissions invalid permission UUID error", func(t *testing.T) {
		mockUow := new(authrepository.MockUnitOfWork)

		permissionUUIDStr := []string{"invalid-uuid"}

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, roleID.String(), permissionUUIDStr)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidRequestBody, err.(*serviceerror.ServiceError).GetErrorMessage())
	})

	t.Run("SyncPermissions FilterValidPermissions error", func(t *testing.T) {
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRoleRepo)
		mockUow.On("PermissionRepository").Return(mockPermissionRepo)

		var permissionUUIDStr []string
		var permissionUUIDs []uuid.UUID
		for _, permission := range permissions {
			permissionUUIDStr = append(permissionUUIDStr, permission.UUID.String())
			permissionUUIDs = append(permissionUUIDs, permission.UUID)
		}

		mockPermissionRepo.On("FilterValidPermissions", permissionUUIDs).Return([]uint64{}, serviceerror.NewServerError())

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, roleID.String(), permissionUUIDStr)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("SyncPermissions invalid role UUID error", func(t *testing.T) {
		mockUow := new(authrepository.MockUnitOfWork)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUow.On("PermissionRepository").Return(mockPermissionRepo)

		var permissionUUIDStr []string
		var permissionUUIDs []uuid.UUID
		var validPermissionIDs []uint64
		for _, permission := range permissions {
			permissionUUIDStr = append(permissionUUIDStr, permission.UUID.String())
			permissionUUIDs = append(permissionUUIDs, permission.UUID)
			validPermissionIDs = append(validPermissionIDs, permission.ID)
		}

		mockPermissionRepo.On("FilterValidPermissions", permissionUUIDs).Return(validPermissionIDs, nil)

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, "invalid-uuid", permissionUUIDStr)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidRequestBody, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("SyncPermissions GetByUUID error", func(t *testing.T) {
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRoleRepo)
		mockUow.On("PermissionRepository").Return(mockPermissionRepo)

		var permissionUUIDStr []string
		var permissionUUIDs []uuid.UUID
		var validPermissionIDs []uint64
		for _, permission := range permissions {
			permissionUUIDStr = append(permissionUUIDStr, permission.UUID.String())
			permissionUUIDs = append(permissionUUIDs, permission.UUID)
			validPermissionIDs = append(validPermissionIDs, permission.ID)
		}

		mockPermissionRepo.On("FilterValidPermissions", permissionUUIDs).Return(validPermissionIDs, nil)
		mockRoleRepo.On("GetByUUID", roleID).Return(&domain.Role{}, serviceerror.New(serviceerror.RecordNotFound))

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, roleID.String(), permissionUUIDStr)

		require.Error(t, err)
		require.Equal(t, serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockPermissionRepo.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("SyncPermissions SyncPermissions error", func(t *testing.T) {
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("RoleRepository").Return(mockRoleRepo)
		mockUow.On("PermissionRepository").Return(mockPermissionRepo)

		var permissionUUIDStr []string
		var permissionUUIDs []uuid.UUID
		var validPermissionIDs []uint64
		for _, permission := range permissions {
			permissionUUIDStr = append(permissionUUIDStr, permission.UUID.String())
			permissionUUIDs = append(permissionUUIDs, permission.UUID)
			validPermissionIDs = append(validPermissionIDs, permission.ID)
		}

		mockRoleRepo.On("GetByUUID", roleID).Return(role, nil)
		mockPermissionRepo.On("FilterValidPermissions", permissionUUIDs).Return(validPermissionIDs, nil)
		mockRoleRepo.On("SyncPermissions", role.ID, validPermissionIDs).Return(serviceerror.NewServerError())

		service := roleservice.New(nil)
		err := service.SyncPermissions(mockUow, roleID.String(), permissionUUIDStr)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockPermissionRepo.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
	})
}
