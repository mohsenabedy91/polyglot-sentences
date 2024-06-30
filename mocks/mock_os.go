package mocks

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type MockOS struct {
	mock.Mock
}

func (m *MockOS) Getwd() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockStat struct {
	mock.Mock
}

func (m *MockStat) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return nil, args.Error(1)
}
