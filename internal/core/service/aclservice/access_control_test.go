package aclservice_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/client"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/aclservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestACLService_CheckAccess(t *testing.T) {

	userUUID := uuid.New()
	user := &domain.User{
		Base: domain.Base{
			ID:   1,
			UUID: userUUID,
		},
	}

	ctx := context.TODO()

	t.Run("SuperAdmin access", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{domain.RoleKeySuperAdmin}, nil)

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.NoError(t, err)
		require.True(t, hasAccess)
		require.Equal(t, user.ID, userID)

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("Access with required permissions", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("PermissionRepository").Return(mockPermissionRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{}, nil)
		mockPermissionRepo.On("GetUserPermissionKeys", user.ID).Return([]domain.PermissionKeyType{domain.PermissionKeyReadUser}, nil)

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.NoError(t, err)
		require.True(t, hasAccess)
		require.Equal(t, user.ID, userID)

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Access denied due to missing required permissions", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("PermissionRepository").Return(mockPermissionRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{}, nil)
		mockPermissionRepo.On("GetUserPermissionKeys", user.ID).Return([]domain.PermissionKeyType{}, nil)

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.NoError(t, err)
		require.False(t, hasAccess)
		require.Equal(t, uint64(0), userID)

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("Access granted due to PermissionKeyNone", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("PermissionRepository").Return(mockPermissionRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{}, nil)
		mockPermissionRepo.On("GetUserPermissionKeys", user.ID).Return([]domain.PermissionKeyType{}, nil)

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyNone)

		require.NoError(t, err)
		require.True(t, hasAccess)
		require.Equal(t, user.ID, userID)

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})

	t.Run("GetByUUID error", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(nil, serviceerror.NewServerError())

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.Error(t, err)
		require.False(t, hasAccess)
		require.Equal(t, uint64(0), userID)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUserClient.AssertExpectations(t)
	})

	t.Run("GetRoleKeys error", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{}, serviceerror.NewServerError())

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.Error(t, err)
		require.False(t, hasAccess)
		require.Equal(t, uint64(0), userID)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("GetUserPermissionKeys error", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockUserClient := new(client.MockUserClient)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockPermissionRepo := new(authrepository.MockPermissionRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("PermissionRepository").Return(mockPermissionRepo)

		mockUserClient.On("GetByUUID", mock.Anything, userUUID.String()).Return(user, nil)
		mockRoleRepo.On("GetRoleKeys", user.ID).Return([]domain.RoleKeyType{}, nil)
		mockPermissionRepo.On("GetUserPermissionKeys", user.ID).Return([]domain.PermissionKeyType{}, serviceerror.NewServerError())

		service := aclservice.New(mockUserClient)
		hasAccess, userID, err := service.CheckAccess(ctx, mockUOW, userUUID, domain.PermissionKeyReadUser)

		require.Error(t, err)
		require.False(t, hasAccess)
		require.Equal(t, uint64(0), userID)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUserClient.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
		mockPermissionRepo.AssertExpectations(t)
	})
}

func TestACLService_AssignUserRoleToUser(t *testing.T) {
	role := &domain.Role{
		Base: domain.Base{
			ID: 1,
		},
		Key: domain.RoleKeyUser,
	}

	t.Run("Assign role successfully", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockACLRepo := new(authrepository.MockACLRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("ACLRepository").Return(mockACLRepo)

		mockRoleRepo.On("GetRoleUser").Return(*role, nil)
		mockACLRepo.On("AssignRolesToUser", uint64(1), []uint64{role.ID}).Return(nil)

		service := aclservice.New(nil)
		err := service.AssignUserRoleToUser(mockUOW, 1)

		require.NoError(t, err)

		mockRoleRepo.AssertExpectations(t)
		mockACLRepo.AssertExpectations(t)
	})

	t.Run("GetRoleUser error", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)

		mockRoleRepo.On("GetRoleUser").Return(domain.Role{}, serviceerror.NewServerError())

		service := aclservice.New(nil)
		err := service.AssignUserRoleToUser(mockUOW, 1)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("AssignRolesToUser error", func(t *testing.T) {
		mockUOW := new(authrepository.MockUnitOfWork)
		mockRoleRepo := new(authrepository.MockRoleRepository)
		mockACLRepo := new(authrepository.MockACLRepository)
		mockUOW.On("RoleRepository").Return(mockRoleRepo)
		mockUOW.On("ACLRepository").Return(mockACLRepo)

		mockRoleRepo.On("GetRoleUser").Return(*role, nil)
		mockACLRepo.On("AssignRolesToUser", uint64(1), []uint64{role.ID}).Return(serviceerror.NewServerError())

		service := aclservice.New(nil)
		err := service.AssignUserRoleToUser(mockUOW, 1)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockRoleRepo.AssertExpectations(t)
		mockACLRepo.AssertExpectations(t)
	})
}
