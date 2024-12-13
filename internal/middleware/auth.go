package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthConfig struct {
	// Add JWT settings or other auth configurations here
}

func AuthMiddleware(config AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Bearer token validation
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		token := parts[1]
		// TODO: Implement token validation logic
		// This is where you would validate the JWT token

		// For demonstration, we'll just check if the token is not empty
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store user information in context for later use
		c.Locals("user", map[string]interface{}{
			"id": "user-id", // This would come from token validation
		})

		return c.Next()
	}
}
