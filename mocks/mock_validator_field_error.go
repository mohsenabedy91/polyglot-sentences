package mocks

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type MockFieldError struct {
	mock.Mock
}

func (r *MockFieldError) Tag() string {
	return r.Called().String(0)
}

func (r *MockFieldError) ActualTag() string {
	return r.Called().String(0)
}

func (r *MockFieldError) Namespace() string {
	return r.Called().String(0)
}

func (r *MockFieldError) StructNamespace() string {
	return r.Called().String(0)
}

func (r *MockFieldError) Field() string {
	return r.Called().String(0)
}

func (r *MockFieldError) StructField() string {
	return r.Called().String(0)
}

func (r *MockFieldError) Value() interface{} {
	return r.Called().Get(0)
}

func (r *MockFieldError) Param() string {
	return r.Called().String(0)
}

func (r *MockFieldError) Kind() reflect.Kind {
	return r.Called().Get(0).(reflect.Kind)
}

func (r *MockFieldError) Type() reflect.Type {
	return r.Called().Get(0).(reflect.Type)
}

func (r *MockFieldError) Translate(ut ut.Translator) string {
	return r.Called(ut).String(0)
}

func (r *MockFieldError) Error() string {
	return r.Called().String(0)
}
