package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey           string
	expiry             time.Duration
	refreshSecretKey   string
	refreshExpiry      time.Duration
}

func NewJWTManager(secretKey, refreshSecretKey string, expiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:         secretKey,
		expiry:           expiry,
		refreshSecretKey: refreshSecretKey,
		refreshExpiry:    refreshExpiry,
	}
}

func (m *JWTManager) Generate(userID, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (m *JWTManager) GenerateTokenPair(userID, email, role string) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = m.Generate(userID, email, role)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with longer expiry
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err = token.SignedString([]byte(m.refreshSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates a refresh token and generates new access and refresh tokens
func (m *JWTManager) RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error) {
	// Parse and validate refresh token
	claims, err := m.validateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Generate new access token
	accessToken, err = m.Generate(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token with the same claims but new expiry
	refreshClaims := Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	newRefreshToken, err = token.SignedString([]byte(m.refreshSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (m *JWTManager) validateRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.refreshSecretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
