package permissionservice_test

import (
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/permissionservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPermissionService_List(t *testing.T) {
	mockLog := new(logger.MockLogger)

	var permissions []*domain.Permission
	for i := 1; i <= 5; i++ {
		permissions = append(permissions, &domain.Permission{
			Base: domain.Base{
				ID:   uint64(i),
				UUID: uuid.New(),
			},
			Title: helper.StringPtr(faker.Word()),
		})
	}

	t.Run("List success", func(t *testing.T) {
		mockRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("PermissionRepository").Return(mockRepo)

		mockRepo.On("List").Return(permissions, nil)

		service := permissionservice.New(mockLog)
		result, err := service.List(mockUow)

		require.NoError(t, err)
		require.NotNil(t, result)
		if len(permissions) > 0 {
			require.Greater(t, len(result), 0)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("List error", func(t *testing.T) {
		mockRepo := new(authrepository.MockPermissionRepository)
		mockUow := new(authrepository.MockUnitOfWork)
		mockUow.On("PermissionRepository").Return(mockRepo)

		mockRepo.On("List").Return([]*domain.Permission{}, serviceerror.NewServerError())

		service := permissionservice.New(mockLog)
		result, err := service.List(mockUow)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
		require.Equal(t, []*domain.Permission{}, result)

		mockRepo.AssertExpectations(t)

	})
}
