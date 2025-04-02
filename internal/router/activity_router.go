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

	// Add some test API keys (in production, these would come from a database)
	// testKey := &apikey.APIKey{
	// 	ID:        "1",
	// 	Name:      "Test App",
	// 	Key:       "dk_test_12345",
	// 	Status:    apikey.KeyStatusActive,
	// 	ExpiresAt: time.Now().AddDate(1, 0, 0), // Expires in 1 year
	// 	RateLimit: 1000,
	// }
	// apiKeyService.AddKey(testKey)

	// Group activity routes
	app.Use(apiKeyMiddleware.ValidateAPIKey())
	activities := app.Group("/api/v1/activities")

	// Apply middleware
	activities.Use(securityMiddleware.JWT())

	// Routes
	activities.Get("/search", activityHandler.Search)
	activities.Post("/", activityHandler.Create)
	activities.Get("/:id", activityHandler.Get)
	activities.Put("/:id", activityHandler.Update)
	//activities.Delete("/:id", activityHandler.Delete)
}
