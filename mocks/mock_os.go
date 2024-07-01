package mocks

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type MockOS struct {
	mock.Mock
}

func (r *MockOS) Getwd() (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}

type MockStat struct {
	mock.Mock
}

func (r *MockStat) Stat(name string) (os.FileInfo, error) {
	args := r.Called(name)
	return nil, args.Error(1)
}
