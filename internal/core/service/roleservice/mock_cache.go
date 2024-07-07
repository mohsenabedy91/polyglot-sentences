package roleservice

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleCacheService struct {
	mock.Mock
}

func (r *MockRoleCacheService) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	args := r.Called(ctx, key)
	return args.Get(0).(*domain.RoleKeyType), args.Error(1)
}

func (r *MockRoleCacheService) SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error {
	args := r.Called(ctx, items)
	return args.Error(0)
}
