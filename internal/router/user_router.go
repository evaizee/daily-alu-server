package router

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/security/apikey"
	"time"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *api.UserHandler, securityMiddleware *middleware.SecurityMiddleware) {
	// Initialize API key service and middleware
	apiKeyService := apikey.NewAPIKeyService()
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)

	// Add some test API keys (in production, these would come from a database)
	testKey := &apikey.APIKey{
		ID:     "1",
		Name:   "Test App",
		Key:    "dk_test_12345",
		Status: apikey.KeyStatusActive,
		ExpiresAt: time.Now().AddDate(1, 0, 0), // Expires in 1 year
		RateLimit: 1000,
	}
	apiKeyService.AddKey(testKey)

	// Apply global middleware
	app.Use(middleware.CORSConfig())
	app.Use(middleware.RateLimiter())
	 // Add API key validation

	// Public routes
	auth := app.Group("/api/v1/auth")

	auth.Get("/verify-email/:token?", userHandler.VerifyEmail)

	app.Use(apiKeyMiddleware.ValidateAPIKey())
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	auth.Use(securityMiddleware.JWT())
	auth.Post("/refresh-token", userHandler.RefreshToken)  // Add refresh token endpoint before API key middleware
	

	// Protected routes
	users := app.Group("/api/v1/users")
	users.Use(securityMiddleware.JWT())

	// Routes accessible by all authenticated users
	users.Patch("/:id/password", userHandler.UpdatePassword)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	

	// Routes accessible only by admins
	admin := users.Group("/")
	admin.Use(securityMiddleware.RoleAuth("admin"))
	admin.Delete("/:id", userHandler.DeleteUser)
}
