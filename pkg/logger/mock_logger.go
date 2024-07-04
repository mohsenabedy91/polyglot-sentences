package logger

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (r *MockLogger) Init(appName string) {
}

func (r *MockLogger) Debug(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) DebugF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Info(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) InfoF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Warn(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) WarnF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Error(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) ErrorF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Fatal(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) FatalF(template string, args ...interface{}) {
	r.Called(template, args)
}
