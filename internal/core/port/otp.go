package port

import "context"

type OtpService interface {
	Set(ctx context.Context, key string, otp string) error
	Validate(ctx context.Context, key string, otp string) error
	Used(ctx context.Context, key string) error
}
