package requests_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/requests"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"testing"
)

func TestCreateUserRequest_ToUserDomain(t *testing.T) {
	tests := []struct {
		name              string
		createUserRequest requests.CreateUserRequest
		expectedResult    domain.User
	}{
		{
			name: "Complete user creation data",
			createUserRequest: requests.CreateUserRequest{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Avatar:    &multipart.FileHeader{Filename: "avatar.png"},
				Gender:    domain.UserGenderOtherStr,
			},
			expectedResult: domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Gender:    domain.UserGenderOtherStr,
			},
		},
		{
			name: "User creation data with no first name",
			createUserRequest: requests.CreateUserRequest{
				FirstName: nil,
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Avatar:    &multipart.FileHeader{Filename: "avatar.png"},
				Gender:    domain.UserGenderOtherStr,
			},
			expectedResult: domain.User{
				FirstName: nil,
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Gender:    domain.UserGenderOtherStr,
			},
		},
		{
			name: "User creation data with no last name",
			createUserRequest: requests.CreateUserRequest{
				FirstName: helper.StringPtr("John"),
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Avatar:    &multipart.FileHeader{Filename: "avatar.png"},
				Gender:    domain.UserGenderPreferNotToSayStr,
			},
			expectedResult: domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Gender:    domain.UserGenderPreferNotToSayStr,
			},
		},
		{
			name: "User creation data with no first and last name",
			createUserRequest: requests.CreateUserRequest{
				FirstName: nil,
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Avatar:    &multipart.FileHeader{Filename: "avatar.png"},
				Gender:    domain.UserGenderMaleStr,
			},
			expectedResult: domain.User{
				FirstName: nil,
				LastName:  nil,
				Email:     "john.doe@gmail.com",
				Gender:    domain.UserGenderMaleStr,
			},
		},
		{
			name: "User creation data with no avatar",
			createUserRequest: requests.CreateUserRequest{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Avatar:    nil,
				Gender:    domain.UserGenderFemaleStr,
			},
			expectedResult: domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Gender:    domain.UserGenderFemaleStr,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := test.createUserRequest.ToUserDomain()
			require.Equal(t, test.expectedResult, user)
		})
	}
}
