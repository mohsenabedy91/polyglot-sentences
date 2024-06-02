package helper_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"math/rand"
	"testing"
	"time"
)

func GenerateOTP(digits int8) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "0123456789"
	otp := make([]byte, digits)
	for i := range otp {
		otp[i] = charset[r.Intn(len(charset))]
	}
	return string(otp)
}

func BenchmarkGenerateOTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		helper.GenerateOTP(6)
	}
}

func BenchmarkGenerateNewOTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateOTP(6)
	}
}
