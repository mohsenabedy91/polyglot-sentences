package authrepository

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetUserPermissionKeys(userID uint64) ([]domain.PermissionKeyType, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.PermissionKeyType), args.Error(1)
}

func (m *MockPermissionRepository) List() ([]*domain.Permission, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FilterValidPermissions(uuids []uuid.UUID) ([]uint64, error) {
	args := m.Called(uuids)
	return args.Get(0).([]uint64), args.Error(1)
}
