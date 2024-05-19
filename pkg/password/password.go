package password

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string, hashedPass chan<- string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)

	hashedPass <- string(bytes)

	return err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
