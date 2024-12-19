package middleware

import (
	"dailyalu-server/internal/security/apikey"
	"github.com/gofiber/fiber/v2"
)

type APIKeyMiddleware struct {
	apiKeyService *apikey.APIKeyService
}

func NewAPIKeyMiddleware(apiKeyService *apikey.APIKeyService) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		apiKeyService: apiKeyService,
	}
}

// ValidateAPIKey middleware validates the API key in the request header
func (m *APIKeyMiddleware) ValidateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get API key from header
		key := c.Get("X-API-Key")
		if key == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key is required",
			})
		}

		// Get client IP
		clientIP := c.IP()

		// Validate key
		apiKey, err := m.apiKeyService.ValidateKey(c.Context(), key, clientIP)
		if err != nil {
			switch err {
			case apikey.ErrInvalidAPIKey:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid API key",
				})
			case apikey.ErrKeyExpired:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "API key has expired",
				})
			case apikey.ErrKeyRevoked:
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "API key has been revoked",
				})
			case apikey.ErrIPNotAllowed:
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "IP address not allowed",
				})
			case apikey.ErrRateLimitExceeded:
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to validate API key",
				})
			}
		}

		// Store API key in context for later use
		c.Locals("api_key", apiKey)

		return c.Next()
	}
}
