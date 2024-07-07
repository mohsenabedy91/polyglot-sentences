package port

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"time"
)

type OTPCacheService interface {
	Set(ctx context.Context, key string, otp string) error
	Validate(ctx context.Context, key string, otp string) error
	Used(ctx context.Context, key string) error

	SetForgetPassword(ctx context.Context, key string, otp string) error
	ValidateForgetPassword(ctx context.Context, key string, otp string) error
	UsedForgetPassword(ctx context.Context, key string) error
}

type OTPCache interface {
	Set(ctx context.Context, key string, value *domain.OTP, expiration time.Duration) error
	Get(ctx context.Context, key string) (*domain.OTP, error)
}
