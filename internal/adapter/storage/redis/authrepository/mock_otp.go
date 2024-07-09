package authrepository

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockOTPCache struct {
	mock.Mock
}

func (r *MockOTPCache) Set(ctx context.Context, key string, value *domain.OTP, expiration time.Duration) error {
	args := r.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (r *MockOTPCache) Get(ctx context.Context, key string) (*domain.OTP, error) {
	args := r.Called(ctx, key)
	return args.Get(0).(*domain.OTP), args.Error(1)
}
