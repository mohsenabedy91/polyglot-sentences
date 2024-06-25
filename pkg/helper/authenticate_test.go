package helper_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"math"
	"strconv"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		cost          int
		expectedError bool
	}{
		{
			name:          "Successful hashing with default cost",
			password:      "mySecurePassword",
			cost:          bcrypt.DefaultCost,
			expectedError: false,
		},
		{
			name:          "Successful hashing with higher cost",
			password:      "mySecurePassword",
			cost:          bcrypt.DefaultCost + 2,
			expectedError: false,
		},
		{
			name:          "Empty password with default cost",
			password:      "",
			cost:          bcrypt.DefaultCost,
			expectedError: false,
		},
		{
			name:          "Invalid cost",
			password:      "mySecurePassword",
			cost:          bcrypt.MaxCost + 1,
			expectedError: true,
		},
		{
			name:          "Password length more than 72 characters",
			password:      "thisisaverylongpasswordthatexceedsthebcryptpasswordlengthlimitof72charactersandshouldbehandledcorrectly",
			cost:          bcrypt.DefaultCost,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := helper.HashPassword(test.password, test.cost)

			if test.expectedError {
				require.Error(t, err)
				require.Equal(t, "", hash)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, hash)

				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(test.password))
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectedResult bool
	}{
		{
			name:     "Correct password and hash",
			password: "mySecurePassword",
			hashedPassword: func() string {
				h, _ := bcrypt.GenerateFromPassword([]byte("mySecurePassword"), bcrypt.DefaultCost)
				return string(h)
			}(),
			expectedResult: true,
		},
		{
			name:     "Incorrect password and hash",
			password: "wrongPassword",
			hashedPassword: func() string {
				h, _ := bcrypt.GenerateFromPassword([]byte("mySecurePassword"), bcrypt.DefaultCost)
				return string(h)
			}(),
			expectedResult: false,
		},
		{
			name:     "Empty password and valid hash",
			password: "",
			hashedPassword: func() string {
				h, _ := bcrypt.GenerateFromPassword([]byte("mySecurePassword"), bcrypt.DefaultCost)
				return string(h)
			}(),
			expectedResult: false,
		},
		{
			name:           "Empty hash",
			password:       "mySecurePassword",
			hashedPassword: "",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			match := helper.CheckPasswordHash(test.password, test.hashedPassword)
			require.Equal(t, test.expectedResult, match)
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	tests := []struct {
		name   string
		digits int8
	}{
		{
			name:   "Generate OTP with 1 digit",
			digits: 1,
		},
		{
			name:   "Generate OTP with 4 digits",
			digits: 4,
		},
		{
			name:   "Generate OTP with 10 digits",
			digits: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			otp := helper.GenerateOTP(test.digits)
			require.NotEmpty(t, otp)
			require.Equal(t, int(test.digits), len(otp))

			minimum := int(math.Pow(10, float64(test.digits-1)))
			maximum := int(math.Pow(10, float64(test.digits))) - 1

			otpInt, err := strconv.Atoi(otp)
			require.NoError(t, err)

			require.GreaterOrEqual(t, otpInt, minimum)
			require.LessOrEqual(t, otpInt, maximum)
		})
	}
}
