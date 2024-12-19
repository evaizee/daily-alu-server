package apikey

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrInvalidAPIKey     = errors.New("invalid API key")
	ErrKeyExpired        = errors.New("API key has expired")
	ErrKeyRevoked        = errors.New("API key has been revoked")
	ErrIPNotAllowed      = errors.New("IP address not allowed")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type APIKeyService struct {
	// In-memory cache for faster lookups
	cache     map[string]*APIKey
	cacheLock sync.RWMutex
}

func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{
		cache: make(map[string]*APIKey),
	}
}

// ValidateKey validates an API key and updates usage statistics
func (s *APIKeyService) ValidateKey(ctx context.Context, keyString string, clientIP string) (*APIKey, error) {
	s.cacheLock.RLock()
	apiKey, exists := s.cache[keyString]
	s.cacheLock.RUnlock()

	if !exists {
		return nil, ErrInvalidAPIKey
	}

	// Check if key is active
	if apiKey.Status != KeyStatusActive {
		return nil, ErrKeyRevoked
	}

	// Check expiration
	if time.Now().After(apiKey.ExpiresAt) {
		return nil, ErrKeyExpired
	}

	// Validate IP if restrictions are in place
	if len(apiKey.AllowedIPs) > 0 {
		allowed := false
		clientIPAddr := net.ParseIP(clientIP)
		for _, allowedIP := range apiKey.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(allowedIP)
			if err == nil && ipNet.Contains(clientIPAddr) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, ErrIPNotAllowed
		}
	}

	// Update last used timestamp
	apiKey.LastUsedAt = time.Now()

	return apiKey, nil
}

// AddKey adds a new API key to the cache
func (s *APIKeyService) AddKey(apiKey *APIKey) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()
	s.cache[apiKey.Key] = apiKey
}

// RemoveKey removes an API key from the cache
func (s *APIKeyService) RemoveKey(key string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()
	delete(s.cache, key)
}
