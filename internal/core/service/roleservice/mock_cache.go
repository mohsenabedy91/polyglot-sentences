package roleservice

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleCache struct {
	mock.Mock
}

func (m *MockRoleCache) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*domain.RoleKeyType), args.Error(1)
}

func (m *MockRoleCache) SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}
