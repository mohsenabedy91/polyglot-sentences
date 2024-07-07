package authrepository

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/stretchr/testify/mock"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (r *MockUnitOfWork) BeginTx(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}

func (r *MockUnitOfWork) Commit() error {
	args := r.Called()
	return args.Error(0)
}

func (r *MockUnitOfWork) Rollback() error {
	args := r.Called()
	return args.Error(0)
}

func (r *MockUnitOfWork) RoleRepository() port.RoleRepository {
	args := r.Called()
	return args.Get(0).(port.RoleRepository)
}

func (r *MockUnitOfWork) PermissionRepository() port.PermissionRepository {
	args := r.Called()
	return args.Get(0).(port.PermissionRepository)
}

func (r *MockUnitOfWork) ACLRepository() port.ACLRepository {
	args := r.Called()
	return args.Get(0).(port.ACLRepository)
}
