package port

import (
	"context"
	"github.com/google/uuid"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	IsEmailUnique(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	VerifiedEmail(ctx context.Context, email string) error
	MarkWelcomeMessageSent(ctx context.Context, ID uint64) error
	UpdateGoogleID(ctx context.Context, ID uint64, googleID string) error
	UpdateLastLoginTime(ctx context.Context, ID uint64) error
	UpdatePassword(ctx context.Context, ID uint64, password string) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	GetByUUID(uow repository.UnitOfWork, uuidStr string) (*domain.User, error)
	GetByID(uow repository.UnitOfWork, id uint64) (*domain.User, error)
	IsEmailUnique(uow repository.UnitOfWork, email string) error
	GetByEmail(uow repository.UnitOfWork, email string) (*domain.User, error)
	List(uow repository.UnitOfWork) ([]domain.User, error)
	Create(uow repository.UnitOfWork, user domain.User) (*domain.User, error)
	VerifiedEmail(uow repository.UnitOfWork, email string) error
	MarkWelcomeMessageSent(uow repository.UnitOfWork, id uint64) error
	UpdateGoogleID(uow repository.UnitOfWork, id uint64, googleID string) error
	UpdateLastLoginTime(uow repository.UnitOfWork, id uint64) error
	UpdatePassword(uow repository.UnitOfWork, id uint64, password string) error
}
