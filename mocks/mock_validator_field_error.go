package mocks

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type MockFieldError struct {
	mock.Mock
}

func (m *MockFieldError) Tag() string {
	return m.Called().String(0)
}

func (m *MockFieldError) ActualTag() string {
	return m.Called().String(0)
}

func (m *MockFieldError) Namespace() string {
	return m.Called().String(0)
}

func (m *MockFieldError) StructNamespace() string {
	return m.Called().String(0)
}

func (m *MockFieldError) Field() string {
	return m.Called().String(0)
}

func (m *MockFieldError) StructField() string {
	return m.Called().String(0)
}

func (m *MockFieldError) Value() interface{} {
	return m.Called().Get(0)
}

func (m *MockFieldError) Param() string {
	return m.Called().String(0)
}

func (m *MockFieldError) Kind() reflect.Kind {
	return m.Called().Get(0).(reflect.Kind)
}

func (m *MockFieldError) Type() reflect.Type {
	return m.Called().Get(0).(reflect.Type)
}

func (m *MockFieldError) Translate(ut ut.Translator) string {
	return m.Called(ut).String(0)
}

func (m *MockFieldError) Error() string {
	return m.Called().String(0)
}
