package config

import (
	"github.com/stretchr/testify/mock"
)

type MockConfiguration struct {
	mock.Mock
}

func (r *MockConfiguration) LoadConfig(envPath ...string) (Config, error) {
	args := r.Called(envPath)
	return args.Get(0).(Config), args.Error(1)
}

func (r *MockConfiguration) GetConfig(envPath ...string) Config {
	args := r.Called(envPath)
	return args.Get(0).(Config)
}
