package port

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	GetByUUID(uuid uuid.UUID) (*domain.User, error)
	GetByID(id uint64) (*domain.User, error)
	IsEmailUnique(email string) (bool, error)
	GetByEmail(email string) (*domain.User, error)
	List() ([]*domain.User, error)
	Save(user *domain.User) (*domain.User, error)
	VerifiedEmail(email string) error
	MarkWelcomeMessageSent(id uint64) error
	UpdateGoogleID(id uint64, googleID string) error
	UpdateLastLoginTime(id uint64) error
	UpdatePassword(id uint64, password string) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	GetByUUID(uow UserUnitOfWork, uuidStr string) (*domain.User, error)
	GetByID(uow UserUnitOfWork, id uint64) (*domain.User, error)
	IsEmailUnique(uow UserUnitOfWork, email string) error
	GetByEmail(uow UserUnitOfWork, email string) (*domain.User, error)
	List(uow UserUnitOfWork) ([]*domain.User, error)
	Create(uow UserUnitOfWork, user domain.User) (*domain.User, error)
	VerifiedEmail(uow UserUnitOfWork, email string) error
	MarkWelcomeMessageSent(uow UserUnitOfWork, id uint64) error
	UpdateGoogleID(uow UserUnitOfWork, id uint64, googleID string) error
	UpdateLastLoginTime(uow UserUnitOfWork, id uint64) error
	UpdatePassword(uow UserUnitOfWork, id uint64, password string) error
}
