package serviceerror_test

import (
	"errors"
	"testing"

	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// normalizeAttributes converts all numeric values in the map to float64
func normalizeAttributes(attrs map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})
	for k, v := range attrs {
		switch value := v.(type) {
		case int:
			normalized[k] = float64(value)
		default:
			normalized[k] = value
		}
	}
	return normalized
}

func TestConvertToGrpcError(t *testing.T) {
	originalAnypbNew := serviceerror.AnypbNew
	defer func() { serviceerror.AnypbNew = originalAnypbNew }()

	tests := []struct {
		name       string
		serviceErr *serviceerror.ServiceError
		setupMock  func()
	}{
		{
			name: "Basic conversion",
			serviceErr: serviceerror.New(
				"Basic error message",
				map[string]interface{}{"key1": "value1", "key2": 123},
			),
			setupMock: func() {},
		},
		{
			name: "Nil attributes",
			serviceErr: serviceerror.New(
				"Error with nil attributes",
				map[string]interface{}{"key": nil},
			),
			setupMock: func() {},
		},
		{
			name: "Error in anypb.New",
			serviceErr: serviceerror.New(
				"Error in anypb.New",
				map[string]interface{}{"key": func() {}},
			),
			setupMock: func() {
				serviceerror.AnypbNew = func(m proto.Message) (*anypb.Any, error) {
					return nil, errors.New("mock error in anypb.New")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setupMock()
			grpcErr := serviceerror.ConvertToGrpcError(test.serviceErr)
			st, ok := status.FromError(grpcErr)
			require.True(t, ok)
			require.Equal(t, test.serviceErr.Error(), st.Message())

			for _, detail := range st.Details() {
				anyDetail, ok := detail.(*anypb.Any)
				require.True(t, ok)

				var customErrorDetail common.CustomErrorDetail
				err := anyDetail.UnmarshalTo(&customErrorDetail)
				require.NoError(t, err)
				require.Equal(t, test.serviceErr.Error(), customErrorDetail.Message)

				expectedAttrs := normalizeAttributes(test.serviceErr.GetAttributes())

				attrs := make(map[string]interface{})
				for k, v := range customErrorDetail.Attributes.Fields {
					attrs[k] = v.AsInterface()
				}
				actualAttrs := normalizeAttributes(attrs)

				require.Equal(t, expectedAttrs, actualAttrs)
			}
		})
	}
}

func TestExtractFromGrpcError(t *testing.T) {
	tests := []struct {
		name        string
		serviceErr  *serviceerror.ServiceError
		expectedErr bool
	}{
		{
			name: "Basic extraction",
			serviceErr: serviceerror.New(
				"Basic error message",
				map[string]interface{}{"key1": "value1", "key2": 123},
			),
			expectedErr: false,
		},
		{
			name:        "Non-status error",
			serviceErr:  serviceerror.New("Non-status error"),
			expectedErr: true,
		},
		{
			name:        "nil error",
			serviceErr:  nil,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var grpcErr error
			if test.name == "Non-status error" {
				grpcErr = errors.New("not a status error")
			} else {
				if test.name != "nil error" {
					grpcErr = serviceerror.ConvertToGrpcError(test.serviceErr)
				}
			}

			extractedErr := serviceerror.ExtractFromGrpcError(grpcErr)
			if test.expectedErr {
				require.NotEqual(t, test.serviceErr, extractedErr)
			} else {
				var extractedServiceErr *serviceerror.ServiceError
				ok := errors.As(extractedErr, &extractedServiceErr)
				require.True(t, ok)
				require.Equal(t, test.serviceErr.Error(), extractedServiceErr.Error())

				expectedAttrs := normalizeAttributes(test.serviceErr.GetAttributes())
				actualAttrs := normalizeAttributes(extractedServiceErr.GetAttributes())

				require.Equal(t, expectedAttrs, actualAttrs)
			}
		})
	}
}
