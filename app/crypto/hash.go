package crypto

import "golang.org/x/crypto/bcrypt"

const hashCost = 12

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
