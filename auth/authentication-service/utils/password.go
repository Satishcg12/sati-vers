package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSalt(len int) (string, error) {
	salt := make([]byte, len)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	// utf8 encoding
	return hex.EncodeToString(hash[:]), nil

}

func ComparePassword(hashedPassword, password string) bool {
	return hashedPassword == password
}
