package helper

import (
	"strconv"
	"strings"
	"unicode"
)

func ConvertToUpperCase(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) {
			builder.WriteRune(unicode.ToUpper(r))
		} else if unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			builder.WriteRune('_')
		}
	}
	return builder.String()
}

func MakeSQLPlaceholders(n int) []string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	return placeholders
}
