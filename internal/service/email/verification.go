package email

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type VerificationService struct {
	tokenExpiry time.Duration
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		tokenExpiry: 24 * time.Hour, // Token valid for 24 hours
	}
}

// GenerateToken creates a secure random token for email verification
func (s *VerificationService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// IsTokenExpired checks if the token has expired based on creation time
func (s *VerificationService) IsTokenExpired(createdAt time.Time) bool {
	return time.Since(createdAt) > s.tokenExpiry
}

// GenerateVerificationLink creates the full verification URL
func (s *VerificationService) GenerateVerificationLink(baseURL, token string) string {
	return fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)
}
