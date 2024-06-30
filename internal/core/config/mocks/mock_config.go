package mocks

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/stretchr/testify/mock"
)

type MockConfiguration struct {
	mock.Mock
}

func (r *MockConfiguration) LoadConfig(envPath ...string) (config.Config, error) {
	args := r.Called(envPath)
	return args.Get(0).(config.Config), args.Error(1)
}

func (r *MockConfiguration) GetConfig(envPath ...string) config.Config {
	args := r.Called(envPath)
	return args.Get(0).(config.Config)
}
