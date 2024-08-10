package helper_test

import (
	"crypto/sha256"
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

func TestToAESKey(t *testing.T) {
	tests := []struct {
		name     string
		inputKey string
	}{
		{
			name:     "Convert short key to AES key",
			inputKey: "short-key",
		},
		{
			name:     "Convert long key to AES key",
			inputKey: "this-is-a-longer-key-to-test-hash-function",
		},
		{
			name:     "Convert empty key to AES key",
			inputKey: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedHash := sha256.Sum256([]byte(test.inputKey))
			expectedKey := expectedHash[:]

			gotKey := helper.ToAESKey(test.inputKey)

			require.Equal(t, expectedKey, gotKey)
			require.Equal(t, 32, len(gotKey))
		})
	}
}

func TestEncryptSecret(t *testing.T) {
	tests := []struct {
		name        string
		secret      []byte
		key         []byte
		expectError bool
	}{
		{
			name:        "Encrypt with valid key and secret",
			secret:      []byte("this is a secret"),
			key:         helper.ToAESKey("my-secret-key"),
			expectError: false,
		},
		{
			name:        "Encrypt with short key",
			secret:      []byte("this is a secret"),
			key:         []byte("short-key"),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encryptedSecret, err := helper.EncryptSecret(test.secret, test.key)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, encryptedSecret)
			}
		})
	}
}

func TestDecryptSecret(t *testing.T) {
	key := helper.ToAESKey("my-secret-key")
	encryptedSecret, err := helper.EncryptSecret([]byte("this is a secret"), key)
	require.NoError(t, err)

	tests := []struct {
		name            string
		encryptedSecret string
		key             []byte
		expectError     bool
		originalSecret  []byte
	}{
		{
			name:            "Decrypt with correct key",
			key:             key,
			originalSecret:  []byte("this is a secret"),
			expectError:     false,
			encryptedSecret: encryptedSecret,
		},
		{
			name:            "Decrypt with wrong key size",
			key:             []byte("w"),
			originalSecret:  []byte("this is a secret"),
			expectError:     true,
			encryptedSecret: encryptedSecret,
		},
		{
			name:            "Decrypt with malformed ciphertext",
			key:             key,
			expectError:     true,
			encryptedSecret: "short-ciphertext",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decryptedSecret, decryptErr := helper.DecryptSecret(test.encryptedSecret, test.key)

			if test.expectError {
				require.Error(t, decryptErr)
			} else {
				require.NoError(t, decryptErr)
				require.Equal(t, test.originalSecret, decryptedSecret)
			}
		})
	}
}
