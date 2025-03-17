package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type TokenType string

const (
	EmailVerification TokenType = "email"
	PasswordReset     TokenType = "password"
)

// TokenService handles secure token generation and verification for various purposes
type TokenService struct {
	configs map[TokenType]time.Duration
}

// NewTokenService creates a new token service with default configurations
func NewTokenService() *TokenService {
	return &TokenService{
		configs: map[TokenType]time.Duration{
			EmailVerification: 24 * time.Hour, // Email verification tokens last 24 hours
			PasswordReset:    1 * time.Hour,   // Password reset tokens last 1 hour for security
		},
	}
}

// GenerateToken creates a secure random token
func (s *TokenService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// IsTokenExpired checks if the token has expired based on creation time and token type
func (s *TokenService) IsTokenExpired(tokenType TokenType, createdAt time.Time) bool {
	expiry, exists := s.configs[tokenType]
	if !exists {
		expiry = 24 * time.Hour // Default expiry
	}
	return time.Since(createdAt) > expiry
}

// GenerateVerificationLink creates the full verification URL for email verification
func (s *TokenService) GenerateVerificationLink(baseURL, token string) string {
	return fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)
}

// GeneratePasswordResetLink creates the full URL for password reset
func (s *TokenService) GeneratePasswordResetLink(baseURL, token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)
}
