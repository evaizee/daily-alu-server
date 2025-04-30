package router

import (
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/security/apikey"
	//"time"

	"github.com/gofiber/fiber/v2"
)

// SetupChildrenRoutes configures the routes for children
func SetupToolsRoutes(
	app *fiber.App,
	securityMiddleware *middleware.SecurityMiddleware,
) {

	// Initialize API key service and middleware
	apiKeyService := apikey.NewAPIKeyService()
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)


	toolsGroup := app.Group("/api/tools", apiKeyMiddleware.ValidateAPIKey())

	// Routes
	toolsGroup.Get("/health-check", func (c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("success")
	})
}
