package authrepository

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRoleCache struct {
	mock.Mock
}

func (r *MockRoleCache) Set(ctx context.Context, key string, value *domain.RoleKeyType, expiration time.Duration) error {
	args := r.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (r *MockRoleCache) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	args := r.Called(ctx, key)
	return args.Get(0).(*domain.RoleKeyType), args.Error(1)
}
