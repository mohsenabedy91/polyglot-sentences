package requests_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRegister_ToUserDomain(t *testing.T) {
	tests := []struct {
		name           string
		authRequest    requests.AuthRegister
		expectedResult domain.User
	}{
		{
			name: "Complete user registration data",
			authRequest: requests.AuthRegister{
				FirstName:         helper.StringPtr("John"),
				LastName:          helper.StringPtr("Doe"),
				Email:             "john.doe@gmail.com",
				Password:          "password",
				ConfirmedPassword: "password",
			},
			expectedResult: domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Password:  helper.StringPtr("password"),
			},
		},
		{
			name: "User registration data with no first name",
			authRequest: requests.AuthRegister{
				FirstName:         nil,
				LastName:          helper.StringPtr("Doe"),
				Email:             "john.doe@gmail.com",
				Password:          "password",
				ConfirmedPassword: "password",
			},
			expectedResult: domain.User{
				FirstName: nil,
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Password:  helper.StringPtr("password"),
			},
		},
		{
			name: "User registration data with no last name",
			authRequest: requests.AuthRegister{
				FirstName:         helper.StringPtr("John"),
				LastName:          nil,
				Email:             "john.doe@gmail.com",
				Password:          "password",
				ConfirmedPassword: "password",
			},
			expectedResult: domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Password:  helper.StringPtr("password"),
			},
		},
		{
			name: "User registration data with no first and last name",
			authRequest: requests.AuthRegister{
				FirstName:         nil,
				LastName:          nil,
				Email:             "john.doe@gmail.com",
				Password:          "password",
				ConfirmedPassword: "password",
			},
			expectedResult: domain.User{
				FirstName: nil,
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Password:  helper.StringPtr("password"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := test.authRequest.ToUserDomain()
			require.Equal(t, test.expectedResult, user)
		})
	}
}
