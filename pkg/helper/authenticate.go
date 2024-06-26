package helper

import (
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOTP(digits int8) string {
	minimum := int(math.Pow(10, float64(digits-1)))
	maximum := int(math.Pow(10, float64(digits)) - 1)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var otp = r.Intn(maximum-minimum) + minimum

	return strconv.Itoa(otp)
}
