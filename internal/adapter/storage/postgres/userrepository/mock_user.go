package userrepository

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) IsEmailUnique(email string) (bool, error) {
	args := r.Called(email)
	return args.Bool(0), args.Error(1)
}

func (r *MockUserRepository) Save(user *domain.User) (*domain.User, error) {
	args := r.Called(user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *MockUserRepository) GetByUUID(uuid uuid.UUID) (*domain.User, error) {
	args := r.Called(uuid)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *MockUserRepository) GetByID(id uint64) (*domain.User, error) {
	args := r.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := r.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *MockUserRepository) List() ([]*domain.User, error) {
	args := r.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (r *MockUserRepository) VerifiedEmail(email string) error {
	args := r.Called(email)
	return args.Error(0)
}

func (r *MockUserRepository) MarkWelcomeMessageSent(id uint64) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *MockUserRepository) UpdateGoogleID(id uint64, googleID string) error {
	args := r.Called(id, googleID)
	return args.Error(0)
}

func (r *MockUserRepository) UpdateLastLoginTime(id uint64) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *MockUserRepository) UpdatePassword(id uint64, password string) error {
	args := r.Called(id, password)
	return args.Error(0)
}
