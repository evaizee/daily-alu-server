package middleware

import (
	"dailyalu-server/internal/security/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// SecurityConfig holds all security-related configurations
type SecurityConfig struct {
	JWTManager *jwt.JWTManager
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware(config SecurityConfig) *SecurityMiddleware {
	return &SecurityMiddleware{
		jwtManager: config.JWTManager,
	}
}

type SecurityMiddleware struct {
	jwtManager *jwt.JWTManager
}

// JWT middleware for authentication
func (m *SecurityMiddleware) JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Remove 'Bearer ' prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := m.jwtManager.Validate(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store user information in context
		c.Locals("user", claims)
		return c.Next()
	}
}

// RoleAuth middleware for role-based authorization
func (m *SecurityMiddleware) RoleAuth(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Claims)
		
		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// CORS middleware configuration
func CORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001,https://dailyalu.mom,http://localhost:5173,",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-API-Key",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           24 * 60 * 60, // 24 hours
	})
}

// RateLimiter middleware configuration
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,                // max number of requests
		Expiration: 1 * time.Minute,    // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // use IP address as key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	})
}
