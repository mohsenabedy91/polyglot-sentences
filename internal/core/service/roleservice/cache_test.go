package roleservice_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/roleservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/require"
)

func TestRoleCacheService_Get(t *testing.T) {
	var role = domain.Role{
		Base: domain.Base{
			ID:   1,
			UUID: uuid.New(),
		},
		Title: "Admin",
		Key:   domain.RoleKeyAdmin,
	}

	ctx := context.TODO()

	t.Run("Get success", func(t *testing.T) {
		mockCache := new(authrepository.MockRoleCache)
		key := fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, role.UUID.String())
		mockCache.On("Get", ctx, key).Return(&role.Key, nil)

		service := roleservice.NewRoleCacheService(mockCache)
		roleKey, err := service.Get(ctx, role.UUID.String())

		require.NoError(t, err)
		require.Equal(t, roleKey, &role.Key)

		mockCache.AssertExpectations(t)
	})

	t.Run("Get cache error", func(t *testing.T) {
		mockCache := new(authrepository.MockRoleCache)
		key := fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, role.UUID.String())
		var expectedRoleKey *domain.RoleKeyType
		mockCache.On("Get", ctx, key).Return(expectedRoleKey, serviceerror.NewServerError())

		service := roleservice.NewRoleCacheService(mockCache)
		roleKey, err := service.Get(ctx, role.UUID.String())

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
		require.Equal(t, roleKey, expectedRoleKey)

		mockCache.AssertExpectations(t)
	})
}

func TestRoleCacheService_SetBulk(t *testing.T) {
	var roles = []domain.Role{
		{
			Base: domain.Base{
				ID:   1,
				UUID: uuid.New(),
			},
			Title: "Admin",
			Key:   domain.RoleKeyAdmin,
		},
		{
			Base: domain.Base{
				ID:   2,
				UUID: uuid.New(),
			},
			Title: "Super Admin",
			Key:   domain.RoleKeySuperAdmin,
		},
	}

	ctx := context.TODO()

	t.Run("SetBulk success", func(t *testing.T) {
		mockCache := new(authrepository.MockRoleCache)

		items := map[string]domain.RoleKeyType{}
		for _, role := range roles {
			items[role.UUID.String()] = role.Key
		}

		for uuidStr, roleKey := range items {
			key := fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, uuidStr)
			mockCache.On("Set", ctx, key, &roleKey, time.Duration(0)).Return(nil)
		}

		service := roleservice.NewRoleCacheService(mockCache)
		err := service.SetBulk(ctx, items)

		require.NoError(t, err)

		mockCache.AssertExpectations(t)
	})

	t.Run("SetBulk cache error", func(t *testing.T) {
		mockCache := new(authrepository.MockRoleCache)

		var uuidStr string
		uuidStr = roles[0].UUID.String()
		var roleKey domain.RoleKeyType
		roleKey = roles[0].Key

		items := map[string]domain.RoleKeyType{}
		for _, role := range roles {
			items[role.UUID.String()] = role.Key
		}

		key := fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, uuidStr)
		mockCache.On("Set", ctx, key, &roleKey, time.Duration(0)).Return(serviceerror.NewServerError())

		service := roleservice.NewRoleCacheService(mockCache)
		err := service.SetBulk(ctx, items)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockCache.AssertExpectations(t)
	})
}
