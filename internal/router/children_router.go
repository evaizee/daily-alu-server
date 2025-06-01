package router

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/security/apikey"
	//"time"

	"github.com/gofiber/fiber/v2"
)

// SetupChildrenRoutes configures the routes for children
func SetupChildrenRoutes(
	app *fiber.App,
	handler *api.ChildrenHandler,
	securityMiddleware *middleware.SecurityMiddleware,
) {

	// Initialize API key service and middleware
	apiKeyService := apikey.NewAPIKeyService()
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)
	
	childrenGroup := app.Group("/children", apiKeyMiddleware.ValidateAPIKey())
	childrenGroup.Use(securityMiddleware.JWT())
	// Routes
	childrenGroup.Post("/", handler.CreateChild)
	childrenGroup.Get("/", handler.GetChildren)
	childrenGroup.Get("/:id", handler.GetChild)
	childrenGroup.Put("/:id", handler.UpdateChild)
}
