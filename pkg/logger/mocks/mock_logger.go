package mocks

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (r *MockLogger) Init(appName string) {
}

func (r *MockLogger) Debug(category logger.Category, subCategory logger.SubCategory, message string, extra map[logger.ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) DebugF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Info(category logger.Category, subCategory logger.SubCategory, message string, extra map[logger.ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) InfoF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Warn(category logger.Category, subCategory logger.SubCategory, message string, extra map[logger.ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) WarnF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Error(category logger.Category, subCategory logger.SubCategory, message string, extra map[logger.ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) ErrorF(template string, args ...interface{}) {
	r.Called(template, args)
}

func (r *MockLogger) Fatal(category logger.Category, subCategory logger.SubCategory, message string, extra map[logger.ExtraKey]interface{}) {
	r.Called(category, subCategory, message, extra)
}

func (r *MockLogger) FatalF(template string, args ...interface{}) {
	r.Called(template, args)
}
