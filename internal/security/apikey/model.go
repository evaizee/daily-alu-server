package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const (
	KeyStatusActive  = "active"
	KeyStatusRevoked = "revoked"
)

type APIKey struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Key         string    `json:"key" db:"key"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at" db:"last_used_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	RateLimit   int       `json:"rate_limit" db:"rate_limit"`
	AllowedIPs  []string  `json:"allowed_ips" db:"allowed_ips"`
}

// GenerateKey generates a new API key with prefix
func GenerateKey(prefix string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + "_" + hex.EncodeToString(bytes), nil
}
