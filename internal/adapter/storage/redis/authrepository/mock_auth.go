package authrepository

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockAuthCache struct {
	mock.Mock
}

func (r *MockAuthCache) SetTokenState(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := r.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (r *MockAuthCache) GetTokenState(ctx context.Context, key string) (string, error) {
	args := r.Called(ctx, key)
	return args.String(0), args.Error(1)
}
