package validations_test

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/validations"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterValidator(t *testing.T) {
	conf := config.Config{
		OTP: config.OTP{
			Digits: 6,
		},
	}

	registerValidatorErr := validations.RegisterValidator(conf)
	assert.NoError(t, registerValidatorErr)

	v, ok := binding.Validator.Engine().(*validator.Validate)
	require.True(t, ok)

	tests := []struct {
		name          string
		tag           string
		fieldValue    interface{}
		expectedValid bool
	}{
		{
			name:          "RegexAlpha valid",
			tag:           "regex_alpha",
			fieldValue:    "test valid",
			expectedValid: true,
		},
		{
			name:          "RegexAlpha invalid",
			tag:           "regex_alpha",
			fieldValue:    "test 123 invalid",
			expectedValid: false,
		},
		{
			name:          "PasswordComplexity valid",
			tag:           "password_complexity",
			fieldValue:    "Password1!",
			expectedValid: true,
		},
		{
			name:          "PasswordComplexity invalid",
			tag:           "password_complexity",
			fieldValue:    "invalid password",
			expectedValid: false,
		},
		{
			name:          "TokenLength valid",
			tag:           "token_length",
			fieldValue:    "123456",
			expectedValid: true,
		},
		{
			name:          "TokenLength invalid",
			tag:           "token_length",
			fieldValue:    "12345",
			expectedValid: false,
		},
		{
			name:          "RoleTitle valid",
			tag:           "role_title",
			fieldValue:    "Admin1",
			expectedValid: true,
		},
		{
			name:          "RoleTitle invalid",
			tag:           "role_title",
			fieldValue:    "Admin@",
			expectedValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := v.Var(test.fieldValue, test.tag)
			if test.expectedValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestRegexAlpha(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedValid bool
	}{
		{
			name:          "Valid string value",
			input:         "abc",
			expectedValid: true,
		},
		{
			name:          "Valid string value contain space",
			input:         "abc def",
			expectedValid: true,
		},
		{
			name:          "Invalid value contain number",
			input:         "abc 123",
			expectedValid: false,
		},
		{
			name:          "Invalid value contain number",
			input:         123,
			expectedValid: false,
		},
		{
			name:          "Invalid value contain special symbol",
			input:         "abc-def",
			expectedValid: false,
		},
		{
			name:          "Invalid nil string value",
			input:         "",
			expectedValid: false,
		},
		{
			name:          "Invalid nil value",
			input:         nil,
			expectedValid: false,
		},
		{
			name:          "Invalid unexpected value",
			input:         func() {},
			expectedValid: false,
		},
	}

	validate := validator.New()
	registerErr := validate.RegisterValidation("regex_alpha", validations.RegexAlpha)
	require.NoError(t, registerErr)

	for _, test := range tests {
		err := validate.Var(test.input, "regex_alpha")
		if test.expectedValid {
			assert.NoError(t, err, "Input: %s", test.input)
		} else {
			assert.Error(t, err, "Input: %s", test.input)
		}
	}
}

func TestPasswordComplexity(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedValid bool
	}{
		{
			name:          "Valid Password contain special symbol",
			input:         "Password1!",
			expectedValid: true,
		},
		{
			name:          "Valid Password contain space",
			input:         "Password1 ",
			expectedValid: true,
		},
		{
			name:          "Invalid Password without special symbol",
			input:         "password1",
			expectedValid: false,
		},
		{
			name:          "Invalid Password without special symbol and lower case char",
			input:         "PASSWORD1",
			expectedValid: false,
		},
		{
			name:          "Invalid Password without digit/number",
			input:         "Password!",
			expectedValid: false,
		},
		{
			name:          "Invalid Password nil value",
			input:         nil,
			expectedValid: false,
		},
		{
			name:          "Invalid unexpected value",
			input:         func() {},
			expectedValid: false,
		},
	}

	validate := validator.New()
	registerErr := validate.RegisterValidation("password_complexity", validations.PasswordComplexity)
	require.NoError(t, registerErr)

	for _, test := range tests {
		err := validate.Var(test.input, "password_complexity")
		if test.expectedValid {
			assert.NoError(t, err, "Input: %s", test.input)
		} else {
			assert.Error(t, err, "Input: %s", test.input)
		}
	}
}

func TestTokenLength(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		length        int8
		expectedValid bool
	}{
		{
			name:          "Valid Length with number value",
			input:         "123456",
			length:        6,
			expectedValid: true,
		},
		{
			name:          "Valid Length start with zero",
			input:         "000001",
			length:        6,
			expectedValid: true,
		},
		{
			name:          "Valid Length with string value",
			input:         "abcdef",
			length:        6,
			expectedValid: true,
		},
		{
			name:          "Invalid Length less than expected value",
			input:         "12345",
			length:        6,
			expectedValid: false,
		},
		{
			name:          "Invalid Length more than expected value",
			input:         "1234567",
			length:        6,
			expectedValid: false,
		},
		{
			name:          "Invalid Length nil value",
			input:         nil,
			length:        6,
			expectedValid: false,
		},
		{
			name:          "Invalid unexpected value",
			input:         func() {},
			length:        6,
			expectedValid: false,
		},
	}

	validate := validator.New()
	registerErr := validate.RegisterValidation("token_length", func(fl validator.FieldLevel) bool {
		return validations.TokenLength(fl, 6)
	})
	require.NoError(t, registerErr)

	for _, test := range tests {
		err := validate.Var(test.input, "token_length")
		if test.expectedValid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestRoleTitle(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedValid bool
	}{
		{
			name:          "Valid Title",
			input:         "Admin",
			expectedValid: true,
		},
		{
			name:          "Valid Title contain number",
			input:         "Admin1",
			expectedValid: true,
		},
		{
			name:          "Valid Title contain space",
			input:         "Admin Role",
			expectedValid: true,
		},
		{
			name:          "Invalid Title contain hyphen",
			input:         "Admin-Role",
			expectedValid: false,
		},
		{
			name:          "Invalid Title contain symbol",
			input:         "Admin@Role",
			expectedValid: false,
		},
		{
			name:          "Invalid Title with nil",
			input:         nil,
			expectedValid: false,
		},
		{
			name:          "Invalid unexpected value",
			input:         func() {},
			expectedValid: false,
		},
	}

	validate := validator.New()
	registerErr := validate.RegisterValidation("role_title", validations.RoleTitle)
	require.NoError(t, registerErr)

	for _, test := range tests {
		err := validate.Var(test.input, "role_title")
		if test.expectedValid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
