package domain_test

import (
	"database/sql"
	"github.com/go-faker/faker/v4"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name           string
		status         domain.UserStatusType
		expectedResult bool
	}{
		{
			name:           "Active user",
			status:         domain.UserStatusActive,
			expectedResult: true,
		},
		{
			name:           "Inactive user",
			status:         domain.UserStatusInactive,
			expectedResult: false,
		},
		{
			name:           "Unverified user",
			status:         domain.UserStatusUnverified,
			expectedResult: false,
		},
		{
			name:           "Banned user",
			status:         domain.UserStatusBanned,
			expectedResult: false,
		},
		{
			name:           "Unknown user status",
			status:         domain.UserStatusUnknown,
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := domain.User{
				Status: test.status,
			}

			require.Equal(t, test.expectedResult, user.IsActive())
		})
	}
}

func TestUser_UserStatusType_String(t *testing.T) {
	tests := []struct {
		name           string
		status         domain.UserStatusType
		expectedResult string
	}{
		{
			name:           "Active user",
			status:         domain.UserStatusActive,
			expectedResult: domain.UserStatusActiveStr,
		},
		{
			name:           "Inactive user",
			status:         domain.UserStatusInactive,
			expectedResult: domain.UserStatusInactiveStr,
		},
		{
			name:           "Unverified user",
			status:         domain.UserStatusUnverified,
			expectedResult: domain.UserStatusUnverifiedStr,
		},
		{
			name:           "Banned user",
			status:         domain.UserStatusBanned,
			expectedResult: domain.UserStatusBannedStr,
		},
		{
			name:           "Unknown user status",
			status:         domain.UserStatusUnknown,
			expectedResult: domain.UserStatusUnknownStr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedResult, test.status.String())
		})
	}
}

func TestUser_ToUserStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedResult domain.UserStatusType
	}{
		{
			name:           "Active user",
			status:         domain.UserStatusActiveStr,
			expectedResult: domain.UserStatusActive,
		},
		{
			name:           "Inactive user",
			status:         domain.UserStatusInactiveStr,
			expectedResult: domain.UserStatusInactive,
		},
		{
			name:           "Unverified user",
			status:         domain.UserStatusUnverifiedStr,
			expectedResult: domain.UserStatusUnverified,
		},
		{
			name:           "Banned user",
			status:         domain.UserStatusBannedStr,
			expectedResult: domain.UserStatusBanned,
		},
		{
			name:           "Unknown user status",
			status:         domain.UserStatusUnknownStr,
			expectedResult: domain.UserStatusUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedResult, domain.ToUserStatus(test.status))
		})
	}
}

func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name           string
		firstName      *string
		lastName       *string
		expectedResult string
	}{
		{
			name:           "Both first and last name are set",
			firstName:      helper.StringPtr("John"),
			lastName:       helper.StringPtr("Doe"),
			expectedResult: "John Doe",
		},
		{
			name:           "Only first name is set",
			firstName:      helper.StringPtr("John"),
			lastName:       nil,
			expectedResult: "John",
		},
		{
			name:           "Only last name is set",
			firstName:      nil,
			lastName:       helper.StringPtr("Doe"),
			expectedResult: "Doe",
		},
		{
			name:           "Neither first nor last name is set",
			firstName:      nil,
			lastName:       nil,
			expectedResult: "",
		},
		{
			name:           "Both first and last name are empty strings",
			firstName:      helper.StringPtr(""),
			lastName:       helper.StringPtr(""),
			expectedResult: "",
		},
		{
			name:           "First name is set, last name is an empty string",
			firstName:      helper.StringPtr("John"),
			lastName:       helper.StringPtr(""),
			expectedResult: "John",
		},
		{
			name:           "First name is an empty string, last name is set",
			firstName:      helper.StringPtr(""),
			lastName:       helper.StringPtr("Doe"),
			expectedResult: "Doe",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := domain.User{
				FirstName: test.firstName,
				LastName:  test.lastName,
			}

			require.Equal(t, test.expectedResult, user.GetFullName())
		})
	}
}

func TestUser_SetGoogleID(t *testing.T) {
	tests := []struct {
		name           string
		input          sql.NullString
		expectedResult *string
	}{
		{
			name:           "Valid google id",
			input:          sql.NullString{String: "valid_google_id", Valid: true},
			expectedResult: helper.StringPtr("valid_google_id"),
		},
		{
			name:           "Invalid google id",
			input:          sql.NullString{String: "invalid_google_id", Valid: false},
			expectedResult: nil,
		},
		{
			name:           "Empty valid google id",
			input:          sql.NullString{String: "", Valid: true},
			expectedResult: helper.StringPtr(""),
		},
		{
			name:           "Empty invalid google id",
			input:          sql.NullString{String: "", Valid: false},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := &domain.User{}

			user.SetGoogleID(test.input)

			if test.expectedResult != nil {
				require.NotNil(t, user.GoogleID)
				require.Equal(t, *test.expectedResult, *user.GoogleID)
			} else {
				require.Nil(t, user.GoogleID)
			}
		})
	}
}

func TestUser_SetFirstName(t *testing.T) {
	tests := []struct {
		name           string
		input          sql.NullString
		expectedResult *string
	}{
		{
			name:           "Valid first name",
			input:          sql.NullString{String: "John", Valid: true},
			expectedResult: helper.StringPtr("John"),
		},
		{
			name:           "Invalid first name",
			input:          sql.NullString{String: faker.FirstName(), Valid: false},
			expectedResult: nil,
		},
		{
			name:           "Empty valid first name",
			input:          sql.NullString{String: "", Valid: true},
			expectedResult: helper.StringPtr(""),
		},
		{
			name:           "Empty invalid first name",
			input:          sql.NullString{String: "", Valid: false},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := &domain.User{}

			user.SetFirstName(test.input)

			if test.expectedResult != nil {
				require.NotNil(t, user.FirstName)
				require.Equal(t, *test.expectedResult, *user.FirstName)
			} else {
				require.Nil(t, user.FirstName)
			}
		})
	}
}

func TestUser_SetLastName(t *testing.T) {
	tests := []struct {
		name           string
		input          sql.NullString
		expectedResult *string
	}{
		{
			name:           "Valid last name",
			input:          sql.NullString{String: "Doe", Valid: true},
			expectedResult: helper.StringPtr("Doe"),
		},
		{
			name:           "Invalid last name",
			input:          sql.NullString{String: faker.LastName(), Valid: false},
			expectedResult: nil,
		},
		{
			name:           "Empty valid last name",
			input:          sql.NullString{String: "", Valid: true},
			expectedResult: helper.StringPtr(""),
		},
		{
			name:           "Empty invalid last name",
			input:          sql.NullString{String: "", Valid: false},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := &domain.User{}

			user.SetLastName(test.input)

			if test.expectedResult != nil {
				require.NotNil(t, user.LastName)
				require.Equal(t, *test.expectedResult, *user.LastName)
			} else {
				require.Nil(t, user.LastName)
			}
		})
	}
}
