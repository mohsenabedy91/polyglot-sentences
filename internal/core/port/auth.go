package port

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"time"
)

type AuthService interface {
	GenerateToken(userUUIDStr string) (*string, error)
	LogoutToken(ctx context.Context, jti string, exp int64) error
}

type UserClient interface {
	IsEmailUnique(ctx context.Context, email string) error

	Create(ctx context.Context, user domain.User) (*domain.User, error)

	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUUID(ctx context.Context, uuidStr string) (*domain.User, error)

	VerifiedEmail(ctx context.Context, email string) error

	MarkWelcomeMessageSent(ctx context.Context, ID uint64) error
	UpdateGoogleID(ctx context.Context, ID uint64, googleID string) error
	UpdateLastLoginTime(ctx context.Context, ID uint64) error
	UpdatePassword(ctx context.Context, ID uint64, password string) error
}

type AuthCache interface {
	SetTokenState(ctx context.Context, key string, value string, expiration time.Duration) error
	GetTokenState(ctx context.Context, key string) (string, error)
}
