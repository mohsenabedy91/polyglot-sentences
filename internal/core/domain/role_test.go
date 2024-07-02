package domain_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRole_SetKey(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedResult domain.Role
	}{
		{
			name: "key with spaces",
			key:  "SUPER ADMIN",
			expectedResult: domain.Role{
				Key: "SUPER_ADMIN",
			},
		},
		{
			name: "key with hyphens",
			key:  "SUPER-admin",
			expectedResult: domain.Role{
				Key: "SUPER_ADMIN",
			},
		},
		{
			name: "key with underscores",
			key:  "super_ADMIN",
			expectedResult: domain.Role{
				Key: "SUPER_ADMIN",
			},
		},
		{
			name: "empty key",
			key:  "",
			expectedResult: domain.Role{
				Key: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			role := domain.Role{}
			role.SetKey(test.key)

			require.NotNil(t, role)
			require.Equal(t, test.expectedResult.Key, role.Key)
		})
	}
}
