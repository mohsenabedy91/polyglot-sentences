package helper

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func MakeSQLPlaceholders(n uint) []string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	return placeholders
}

func StringPtr(str string) *string {
	return &str
}

func RandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var result = make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = Chars[rand.Intn(len(Chars))]
	}

	return string(result)
}
