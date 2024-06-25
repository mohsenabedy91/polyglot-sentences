package mocks

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/stretchr/testify/mock"
)

type MockConfiguration struct {
	mock.Mock
}

func (m *MockConfiguration) LoadConfig(envPath ...string) (config.Config, error) {
	args := m.Called(envPath)
	return args.Get(0).(config.Config), args.Error(1)
}

func (m *MockConfiguration) GetConfig(envPath ...string) config.Config {
	args := m.Called(envPath)
	return args.Get(0).(config.Config)
}
