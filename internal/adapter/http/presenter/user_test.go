package presenter_test

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/proto/user"
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
				Gender:    domain.UserGenderMale,
				AuthChallenges: []*user.AuthChallenge{
					{
						Type:   domain.ChallengeTOTPStr,
						Status: domain.ChallengeEnableStr,
					},
				},
			},
			expectedResult: &presenter.User{
				ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
				FirstName: helper.StringPtr("John"),
				LastName:  helper.StringPtr("Doe"),
				Email:     "john.doe@gmail.com",
				Status:    string(domain.UserStatusActive),
				Gender:    helper.StringPtr("MALE"),
				AuthChallenges: []presenter.AuthChallenge{
					{
						Type:   "TOTP",
						Status: "ENABLE",
					},
				},
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
		{
			name: "Valid User with Gender",
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
				Gender:    helper.StringPtr(domain.UserGenderPreferNotToSayStr),
			},
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
				Gender:    helper.StringPtr(domain.UserGenderPreferNotToSayStr),
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
			name:           "Empty user list",
			users:          []*domain.User{},
			expectedResult: nil,
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
					Gender:    domain.UserGenderPreferNotToSay,
				},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("fbed3952-feac-4165-958c-202d8a1c80b7"),
					},
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    domain.UserStatusInactive,
					Gender:    domain.UserGenderPreferNotToSay,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    string(domain.UserStatusActive),
					Gender:    helper.StringPtr(domain.UserGenderPreferNotToSayStr),
				},
				{
					ID:        "fbed3952-feac-4165-958c-202d8a1c80b7",
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    string(domain.UserStatusInactive),
					Gender:    helper.StringPtr(domain.UserGenderPreferNotToSayStr),
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
					FirstName:      helper.StringPtr("John"),
					LastName:       helper.StringPtr("Doe"),
					Email:          "john.doe@gmail.com",
					Status:         domain.UserStatusActive,
					Gender:         domain.UserGenderPreferNotToSay,
					AuthChallenges: nil,
				},
				{},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("fbed3952-feac-4165-958c-202d8a1c80b7"),
					},
					FirstName:      helper.StringPtr("Jane"),
					LastName:       helper.StringPtr("Smith"),
					Email:          "jane.smith@gmail.com",
					Status:         domain.UserStatusInactive,
					Gender:         domain.UserGenderPreferNotToSay,
					AuthChallenges: nil,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:             "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName:      helper.StringPtr("John"),
					LastName:       helper.StringPtr("Doe"),
					Email:          "john.doe@gmail.com",
					Status:         string(domain.UserStatusActive),
					Gender:         helper.StringPtr(domain.UserGenderPreferNotToSayStr),
					AuthChallenges: nil,
				},
				{
					ID:             "fbed3952-feac-4165-958c-202d8a1c80b7",
					FirstName:      helper.StringPtr("Jane"),
					LastName:       helper.StringPtr("Smith"),
					Email:          "jane.smith@gmail.com",
					Status:         string(domain.UserStatusInactive),
					Gender:         helper.StringPtr(domain.UserGenderPreferNotToSayStr),
					AuthChallenges: nil,
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
					Gender:    domain.UserGenderPreferNotToSay,
				},
				{},
				{
					FirstName: helper.StringPtr("Jane"),
					LastName:  helper.StringPtr("Smith"),
					Email:     "jane.smith@gmail.com",
					Status:    domain.UserStatusInactive,
					Gender:    domain.UserGenderPreferNotToSay,
				},
			},
			expectedResult: []presenter.User{
				{
					ID:        "2b1ef850-5b3a-441e-bd26-33f50e527b7a",
					FirstName: helper.StringPtr("John"),
					LastName:  helper.StringPtr("Doe"),
					Email:     "john.doe@gmail.com",
					Status:    string(domain.UserStatusActive),
					Gender:    helper.StringPtr(domain.UserGenderPreferNotToSayStr),
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

func TestToTOTPResource(t *testing.T) {
	tests := []struct {
		name           string
		totp           *domain.TOTPKey
		expectedResult *presenter.TOTPKey
	}{
		{
			name:           "Nil TOTP",
			totp:           nil,
			expectedResult: nil,
		},
		{
			name: "Valid TOTP",
			totp: &domain.TOTPKey{
				Secret: "this is a secret",
				URL:    "otpauth_url",
			},
			expectedResult: &presenter.TOTPKey{
				Secret: "this is a secret",
				URL:    "otpauth_url",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToTOTPResource(test.totp)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestPrepareAuthChallenge(t *testing.T) {
	tests := []struct {
		name           string
		authChallenges []*user.AuthChallenge
		expectedResult []presenter.AuthChallenge
	}{
		{
			name:           "Empty authChallenges",
			authChallenges: []*user.AuthChallenge{},
			expectedResult: nil,
		},
		{
			name: "Valid authChallenges",
			authChallenges: []*user.AuthChallenge{
				{
					Type:   domain.ChallengeTOTPStr,
					Status: domain.ChallengeEnableStr,
				},
				{
					Type:   domain.ChallengeRecoveryCodesStr,
					Status: domain.ChallengeDisableStr,
				},
			},
			expectedResult: []presenter.AuthChallenge{
				{
					Type:   "TOTP",
					Status: "ENABLE",
				},
				{
					Type:   "RECOVERY_CODES",
					Status: "DISABLE",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.PrepareAuthChallenge(test.authChallenges)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
