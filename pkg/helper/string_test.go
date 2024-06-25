package helper_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertToUpperCase(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult string
	}{
		{
			name:           "Lowercase with spaces",
			input:          "all characters are lower",
			expectedResult: "ALL_CHARACTERS_ARE_LOWER",
		},
		{
			name:           "Lowercase with hyphens",
			input:          "all-characters-are-lower",
			expectedResult: "ALL_CHARACTERS_ARE_LOWER",
		},
		{
			name:           "Lowercase with underscores",
			input:          "all_characters_are_lower",
			expectedResult: "ALL_CHARACTERS_ARE_LOWER",
		},
		{
			name:           "Mixed with numbers",
			input:          "there is number 66",
			expectedResult: "THERE_IS_NUMBER_66",
		},
		{
			name:           "Uppercase with spaces",
			input:          "ALL CHARACTERS ARE UPPER CASE",
			expectedResult: "ALL_CHARACTERS_ARE_UPPER_CASE",
		},
		{
			name:           "Exact match",
			input:          "SAME_EXPECTED_RESULT",
			expectedResult: "SAME_EXPECTED_RESULT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helper.ConvertToUpperCase(test.input)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestMakeSQLPlaceholders(t *testing.T) {
	tests := []struct {
		name           string
		n              uint
		expectedResult []string
	}{
		{
			name:           "Zero placeholders",
			n:              0,
			expectedResult: []string{},
		},
		{
			name:           "One placeholder",
			n:              1,
			expectedResult: []string{"$1"},
		},
		{
			name:           "Two placeholders",
			n:              2,
			expectedResult: []string{"$1", "$2"},
		},
		{
			name:           "Multiple placeholders",
			n:              5,
			expectedResult: []string{"$1", "$2", "$3", "$4", "$5"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helper.MakeSQLPlaceholders(test.n)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
