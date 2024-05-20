package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	IsEmailUnique(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *domain.User) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	IsEmailUnique(ctx context.Context, email string) error
	RegisterUser(ctx context.Context, user domain.User) error

	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}
