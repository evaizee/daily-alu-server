package utils

import (
	"dailyalu-server/internal/security/jwt"
	"github.com/gofiber/fiber/v2"
)

// GetUserFromContext extracts the user claims from the Fiber context
func GetUserFromContext(c *fiber.Ctx) *jwt.Claims {
	return c.Locals("user").(*jwt.Claims)
}

// GetUserIDFromContext extracts just the user ID from the context
func GetUserIDFromContext(c *fiber.Ctx) string {
	claims := GetUserFromContext(c)
	return claims.UserID
}

// IsAdminUser checks if the user in context has admin role
func IsAdminUser(c *fiber.Ctx) bool {
	claims := GetUserFromContext(c)
	return claims.Role == "admin"
}

// CanAccessUserData checks if the requesting user can access data for the specified user ID
// (either they are that user or they're an admin)
func CanAccessUserData(c *fiber.Ctx, userID string) bool {
	claims := GetUserFromContext(c)
	return claims.UserID == userID || claims.Role == "admin"
}
