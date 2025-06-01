package router

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/security/apikey"
	//"time"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *api.UserHandler, securityMiddleware *middleware.SecurityMiddleware) {
	// Initialize API key service and middleware
	apiKeyService := apikey.NewAPIKeyService()
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)

	// Apply global middleware
	app.Use(middleware.CORSConfig())
	
	// Add API key validation
	app.Use(apiKeyMiddleware.ValidateAPIKey())

	// Public routes
	auth := app.Group("/v1/auth")

	auth.Get("/verify-email/:token", userHandler.VerifyEmail)

	// Apply rate limiters to login and register endpoints
	middleware.RateLimitedRoute(auth, "POST", "/register", userHandler.Register)
	middleware.RateLimitedRoute(auth, "POST", "/login", userHandler.Login)
	
	// Password recovery routes (don't require authentication)
	auth.Post("/forgot-password", userHandler.ForgotPassword)
	auth.Post("/reset-password", userHandler.ResetPassword)

	auth.Use(securityMiddleware.JWT())
	auth.Post("/refresh-token", userHandler.RefreshToken)  // Add refresh token endpoint before API key middleware

	// Protected routes
	users := app.Group("/api/v1/users")
	users.Use(securityMiddleware.JWT())

	// Routes accessible by all authenticated users
	users.Patch("/password", userHandler.UpdatePassword)
	users.Get("/profile", userHandler.GetUser)
	users.Put("/profile", userHandler.UpdateUser)
	

	// Routes accessible only by admins
	admin := users.Group("/")
	admin.Use(securityMiddleware.RoleAuth("admin"))
	admin.Delete("/:id", userHandler.DeleteUser)
	admin.Get("/:id", userHandler.AdminGetUser)  // Admin-specific route to get any user
	admin.Put("/:id", userHandler.AdminUpdateUser)  // Admin-specific route to update any user
}
