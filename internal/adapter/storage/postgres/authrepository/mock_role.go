package authrepository

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (r *MockRoleRepository) Create(role domain.Role) error {
	args := r.Called(role)
	return args.Error(0)
}

func (r *MockRoleRepository) GetByUUID(uuid uuid.UUID) (*domain.Role, error) {
	args := r.Called(uuid)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (r *MockRoleRepository) List() ([]*domain.Role, error) {
	args := r.Called()
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (r *MockRoleRepository) Update(role domain.Role, uuid uuid.UUID) error {
	args := r.Called(role, uuid)
	return args.Error(0)
}

func (r *MockRoleRepository) Delete(uuid uuid.UUID, deletedBy uint64) error {
	args := r.Called(uuid, deletedBy)
	return args.Error(0)
}

func (r *MockRoleRepository) ExistKey(key domain.RoleKeyType) (bool, error) {
	args := r.Called(key)
	return args.Bool(0), args.Error(1)
}

func (r *MockRoleRepository) GetRoleUser() (domain.Role, error) {
	args := r.Called()
	return args.Get(0).(domain.Role), args.Error(1)
}

func (r *MockRoleRepository) GetPermissions(roleUUID uuid.UUID) (*domain.Role, error) {
	args := r.Called(roleUUID)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (r *MockRoleRepository) SyncPermissions(roleID uint64, permissionIDs []uint64) error {
	args := r.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (r *MockRoleRepository) GetUserRoleKeys(userID uint64) ([]domain.RoleKeyType, error) {
	args := r.Called(userID)
	return args.Get(0).([]domain.RoleKeyType), args.Error(1)
}
