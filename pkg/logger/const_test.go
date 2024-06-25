package logger_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapToZapParams(t *testing.T) {
	tests := []struct {
		name           string
		input          map[logger.ExtraKey]interface{}
		expectedResult []interface{}
	}{
		{
			name:           "Empty map",
			input:          map[logger.ExtraKey]interface{}{},
			expectedResult: []interface{}{},
		},
		{
			name: "Map with one key-value pair string type",
			input: map[logger.ExtraKey]interface{}{
				"key1": "value1",
			},
			expectedResult: []interface{}{"key1", "value1"},
		},
		{
			name: "Map with one key-value pair integer type",
			input: map[logger.ExtraKey]interface{}{
				"key1": 1,
			},
			expectedResult: []interface{}{"key1", 1},
		},
		{
			name: "Map with multiple key-value pairs and types",
			input: map[logger.ExtraKey]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": true,
			},
			expectedResult: []interface{}{"key1", "value1", "key2", 2, "key3", true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := logger.MapToZapParams(test.input)
			require.ElementsMatch(t, test.expectedResult, result)
		})
	}
}
