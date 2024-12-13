package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	minCost     = 10
	maxCost     = 12
	defaultCost = minCost
)

// Hash creates a bcrypt hash of the password using the default cost
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify checks if the password matches the hash
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
