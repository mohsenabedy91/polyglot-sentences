package userrepository

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

func (r *MockUnitOfWork) UserRepository() port.UserRepository {
	args := r.Called()
	return args.Get(0).(port.UserRepository)
}
