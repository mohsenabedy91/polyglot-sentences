package port

import "context"

type OtpService interface {
	Set(ctx context.Context, key string, otp string) error
	Validate(ctx context.Context, key string, otp string) error
	Used(ctx context.Context, key string) error

	SetForgetPassword(ctx context.Context, key string, otp string) error
	ValidateForgetPassword(ctx context.Context, key string, otp string) error
	UsedForgetPassword(ctx context.Context, key string) error
}
