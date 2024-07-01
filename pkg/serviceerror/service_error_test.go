package serviceerror_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

func TestNewServerError(t *testing.T) {
	err := serviceerror.NewServerError()
	require.NotNil(t, err)

	require.Equal(t, err.Error(), string(serviceerror.ServerError))
	require.Equal(t, err.GetErrorMessage(), serviceerror.ServerError)
	require.Nil(t, err.GetAttributes())
}

func TestServiceErrorMethods(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr *serviceerror.ServiceError
	}{
		{
			name:       "Server error",
			serviceErr: serviceerror.New(serviceerror.ServerError),
		},
		{
			name: "Email registered error",
			serviceErr: serviceerror.New(
				serviceerror.EmailRegistered,
				map[string]interface{}{"key1": "value1", "key2": 123},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Server error" {
				require.Equal(t, test.serviceErr.GetErrorMessage(), serviceerror.ServerError)
				require.Equal(t, test.serviceErr.GetAttributes(), map[string]interface{}(nil))
				require.Equal(t, test.serviceErr.Error(), string(serviceerror.ServerError))
			} else {
				require.Equal(t, test.serviceErr.GetErrorMessage(), serviceerror.EmailRegistered)
				require.Equal(t, test.serviceErr.GetAttributes(), map[string]interface{}{"key1": "value1", "key2": 123})
				require.Equal(t, test.serviceErr.Error(), string(serviceerror.EmailRegistered))
			}
		})
	}
}
