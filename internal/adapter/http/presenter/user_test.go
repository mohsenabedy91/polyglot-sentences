package presenter_test

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrepareUser(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		expectedResult *presenter.User
	}{
		{
			name:           "Nil user",
			user:           nil,
			expectedResult: nil,
		},
		{
			name: "Valid User",
			user: &domain.User{
				Base: domain.Base{
					UUID: uuid.MustParse("2b1ef850-5b3a-441e-bd26-33f50e527b7a"),
				},
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    domain.UserStatusActive,
			},
			expectedResult: &presenter.User{
				ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    string(domain.UserStatusActive),
			},
		},
		{
			name: "Invalid user with uuid equal nil",
			user: &domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    domain.UserStatusActive,
			},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.PrepareUser(test.user)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToUserResource(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		expectedResult *presenter.User
	}{
		{
			name:           "Nil user",
			user:           nil,
			expectedResult: nil,
		},
		{
			name: "Valid user",
			user: &domain.User{
				Base: domain.Base{
					UUID: uuid.MustParse("2b1ef850-5b3a-441e-bd26-33f50e527b7a"),
				},
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    domain.UserStatusActive,
			},
			expectedResult: &presenter.User{
				ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    string(domain.UserStatusActive),
			},
		},
		{
			name: "Invalid user with uuid equal nil",
			user: &domain.User{
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    domain.UserStatusActive,
			},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToUserResource(test.user)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToUserCollection(t *testing.T) {
	tests := []struct {
		name           string
		users          []*domain.User
		expectedResult []presenter.User
	}{
		{
			name: "Valid Users",
			users: []*domain.User{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("2b1ef850-5b3a-441e-bd26-33f50e527b7a"),
					},
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    domain.UserStatusActive,
				},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("fbed3952-feac-4165-958c-202d8a1c80b7"),
					},
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    domain.UserStatusInactive,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    string(domain.UserStatusActive),
				},
				{
					ID:        "fbed3952-feac-4165-958c-202d8a1c80b7",
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    string(domain.UserStatusInactive),
				},
			},
		},
		{
			name: "Valid Users",
			users: []*domain.User{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("2b1ef850-5b3a-441e-bd26-33f50e527b7a"),
					},
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    domain.UserStatusActive,
				},
				{},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("fbed3952-feac-4165-958c-202d8a1c80b7"),
					},
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    domain.UserStatusInactive,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    string(domain.UserStatusActive),
				},
				{
					ID:        "fbed3952-feac-4165-958c-202d8a1c80b7",
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    string(domain.UserStatusInactive),
				},
			},
		},
		{
			name: "Invalid user with uuid equal nil",
			users: []*domain.User{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("2b1ef850-5b3a-441e-bd26-33f50e527b7a"),
					},
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    domain.UserStatusActive,
				},
				{},
				{
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    domain.UserStatusInactive,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    string(domain.UserStatusActive),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToUserCollection(test.users)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
