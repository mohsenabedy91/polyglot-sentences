package claim_test

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/claim"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserUUIDFromGinContext(t *testing.T) {
	tests := []struct {
		name             string
		userUUIDStr      string
		expectedResult   uuid.UUID
		expectedError    bool
		shouldSetContext bool
	}{
		{
			name:             "Valid UUID in context",
			userUUIDStr:      "123e4567-e89b-12d3-a456-426614174000",
			expectedResult:   uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedError:    false,
			shouldSetContext: true,
		},
		{
			name:             "Invalid UUID in context",
			userUUIDStr:      "invalid-uuid",
			expectedError:    true,
			shouldSetContext: true,
		},
		{
			name:             "Missing UUID in context",
			userUUIDStr:      "",
			expectedError:    true,
			shouldSetContext: true,
		},
		{
			name:             "Missing set context",
			expectedError:    true,
			shouldSetContext: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := &gin.Context{}
			if test.shouldSetContext {
				ctx.Set(config.AuthTokenUserUUID, test.userUUIDStr)
			}

			if test.expectedError {
				require.Panics(t, func() {
					claim.GetUserUUIDFromGinContext(ctx)
				})
			} else {
				result := claim.GetUserUUIDFromGinContext(ctx)
				require.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestGetJTIFromGinContext(t *testing.T) {
	tests := []struct {
		name             string
		jti              string
		expectedResult   string
		shouldSetContext bool
	}{
		{
			name:             "Valid JTI in context",
			jti:              "some-jti-value",
			expectedResult:   "some-jti-value",
			shouldSetContext: true,
		},
		{
			name:             "Missing JTI in context",
			jti:              "",
			expectedResult:   "",
			shouldSetContext: true,
		},
		{
			name:             "Missing set context",
			shouldSetContext: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := &gin.Context{}
			if test.shouldSetContext {
				ctx.Set(config.AuthTokenJTI, test.jti)
			}

			result := claim.GetJTIFromGinContext(ctx)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestGetExpFromGinContext(t *testing.T) {
	tests := []struct {
		name             string
		expirationTime   float64
		expectedResult   int64
		shouldSetContext bool
	}{
		{
			name:             "Valid expiration time in context",
			expirationTime:   1627840201.0,
			expectedResult:   1627840201,
			shouldSetContext: true,
		},
		{
			name:             "Missing expiration time in context",
			expirationTime:   0,
			expectedResult:   0,
			shouldSetContext: true,
		},
		{
			name:             "Missing set context",
			shouldSetContext: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := &gin.Context{}
			if test.shouldSetContext {
				ctx.Set(config.AuthTokenExpirationTime, test.expirationTime)
			}

			result := claim.GetExpFromGinContext(ctx)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
