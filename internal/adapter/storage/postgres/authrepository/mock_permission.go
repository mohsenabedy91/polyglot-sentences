package authrepository

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (r *MockPermissionRepository) GetUserPermissionKeys(userID uint64) ([]domain.PermissionKeyType, error) {
	args := r.Called(userID)
	return args.Get(0).([]domain.PermissionKeyType), args.Error(1)
}

func (r *MockPermissionRepository) List() ([]*domain.Permission, error) {
	args := r.Called()
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

func (r *MockPermissionRepository) FilterValidPermissions(uuids []uuid.UUID) ([]uint64, error) {
	args := r.Called(uuids)
	return args.Get(0).([]uint64), args.Error(1)
}
