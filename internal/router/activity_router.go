package router

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/security/apikey"
	//"time"

	"github.com/gofiber/fiber/v2"
)

func SetupActivityRoutes(app *fiber.App, activityHandler *api.ActivityHandler, securityMiddleware *middleware.SecurityMiddleware) {
	// Initialize API key service and middleware
	apiKeyService := apikey.NewAPIKeyService()
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)

	// Group activity routes
	app.Use(apiKeyMiddleware.ValidateAPIKey())
	activities := app.Group("/v1/activities")

	// Apply middleware
	activities.Use(securityMiddleware.JWT())

	// Routes
	activities.Get("/search", activityHandler.Search)
	activities.Post("/", activityHandler.Create)
	activities.Get("/:id", activityHandler.Get)
	activities.Put("/:id", activityHandler.Update)
	//activities.Delete("/:id", activityHandler.Delete)
}
