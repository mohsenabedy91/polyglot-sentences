package helper

import (
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
