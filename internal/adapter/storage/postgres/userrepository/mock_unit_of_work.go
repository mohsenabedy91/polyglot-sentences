package userrepository

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/stretchr/testify/mock"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) BeginTx(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUnitOfWork) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUnitOfWork) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUnitOfWork) UserRepository() port.UserRepository {
	args := m.Called()
	return args.Get(0).(port.UserRepository)
}
