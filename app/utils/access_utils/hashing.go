package access_utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Generate hashed string from plaintext
func GeneratePasswordHash(plainPassword string) (string, error) {
	password := []byte(plainPassword)

	bytes, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		fmt.Println("[utils]", err.Error())
		return "", err
	}

	hashedPassword := string(bytes)
	return hashedPassword, nil
}

// Verify if input password matches the encrypted password saved in DB
func VerifyPasswordHash(loginPassword string, storedPassword string) bool {
	password := []byte(loginPassword)
	hash := []byte(storedPassword)

	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		fmt.Println("[utils]", err.Error())
		return false
	}
	return true
}
