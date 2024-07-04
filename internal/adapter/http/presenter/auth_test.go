package presenter_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToTokenResource(t *testing.T) {
	tests := []struct {
		name           string
		token          *string
		expectedResult *presenter.Token
	}{
		{
			name:  "Non-nil token",
			token: helper.StringPtr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			expectedResult: &presenter.Token{
				AccessToken: helper.StringPtr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			},
		},
		{
			name:  "Nil token",
			token: nil,
			expectedResult: &presenter.Token{
				AccessToken: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToTokenResource(test.token)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
