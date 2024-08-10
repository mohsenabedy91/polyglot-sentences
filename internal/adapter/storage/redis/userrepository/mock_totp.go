package userrepository

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockTOTPCache struct {
	mock.Mock
}

func (r *MockTOTPCache) Set(ctx context.Context, key string, value string) error {
	args := r.Called(ctx, key, value)
	return args.Error(0)
}

func (r *MockTOTPCache) Get(ctx context.Context, key string) (string, error) {
	args := r.Called(ctx, key)
	return args.String(0), args.Error(1)
}
