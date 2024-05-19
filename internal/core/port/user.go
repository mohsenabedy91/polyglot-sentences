package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	IsEmailUnique(ctx context.Context, email string) (bool, serviceerror.Error)
	Save(ctx context.Context, user *domain.User) serviceerror.Error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, serviceerror.Error)
	GetByEmail(ctx context.Context, email string) (*domain.User, serviceerror.Error)
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	IsEmailUnique(ctx context.Context, email string) serviceerror.Error
	RegisterUser(ctx context.Context, user domain.User) serviceerror.Error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, serviceerror.Error)
	GetByEmail(ctx context.Context, email string) (*domain.User, serviceerror.Error)
}
