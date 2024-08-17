package helper

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
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

// ToAESKey Convert key to 32 bytes using SHA-256
func ToAESKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptSecret encrypts the TOTP secret using AES-256.
func EncryptSecret(secret, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(secret))
	iv := cipherText[:aes.BlockSize]
	if _, ioErr := io.ReadFull(cryptorand.Reader, iv); ioErr != nil {
		return "", ioErr
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], secret)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// DecryptSecret decrypts the encrypted TOTP secret.
func DecryptSecret(encryptedSecret string, key []byte) ([]byte, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(encryptedSecret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}
