package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Algon2Params struct {
	Memory      uint32 // memory in KB (it represents the amount of memory that the algorithm uses)
	Iterations  uint32 // number of iterations (it represents the number of iterations that the algorithm performs)
	Parallelism uint8  // parallelism (it represents the number of threads that the algorithm uses)
	SaltLength  uint32 // salt length (it represents the length of the salt)
	KeyLength   uint32 // key length (it represents the length of the key)
}

func GenerateString(len int) (string, error) {
	salt := make([]byte, len)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func HashPasswordWithArgon2(password string, params *Algon2Params) (string, error) {
	// generate salt
	salt, err := GenerateString(int(params.SaltLength))
	if err != nil {
		return "", err
	}
	// hash password
	hash := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
	b64Salt := hex.EncodeToString([]byte(salt))
	b64Hash := hex.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func VerifyPasswordWithArgon2(password, encodedHash string) (bool, error) {
	// decode hash
	var version, memory, iterations, parallelism int
	var salt, hash []byte
	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", &version, &memory, &iterations, &parallelism, &salt, &hash)
	if err != nil {
		return false, err
	}
	// hash password
	decodedHash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(iterations),
		uint32(memory),
		uint8(parallelism),
		uint32(len(hash)),
	)
	return subtle.ConstantTimeCompare(hash, decodedHash) == 1, nil
}

func HashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	// utf8 encoding
	return hex.EncodeToString(hash[:]), nil

}

func VerifyPassword(password, salt, hashedPassword string) bool {
	tempHashedPassword, err := HashPassword(password, salt)
	if err != nil {
		return false
	}
	return tempHashedPassword == hashedPassword
}
