package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	IsEmailUnique(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	GetByUUID(ctx context.Context, uuidStr string) (*domain.User, error)
	IsEmailUnique(ctx context.Context, email string) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) error
}
