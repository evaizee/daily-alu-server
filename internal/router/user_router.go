package router

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *api.UserHandler, securityMiddleware *middleware.SecurityMiddleware) {
	// Apply global middleware
	app.Use(middleware.CORSConfig())
	app.Use(middleware.RateLimiter())

	// Public routes
	auth := app.Group("/api/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Protected routes
	users := app.Group("/api/users")
	users.Use(securityMiddleware.JWT())

	// Routes accessible by all authenticated users
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)

	// Routes accessible only by admins
	admin := users.Group("/")
	admin.Use(securityMiddleware.RoleAuth("admin"))
	admin.Delete("/:id", userHandler.DeleteUser)
}
