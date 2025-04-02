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

	// Add some test API keys (in production, these would come from a database)
	// testKey := &apikey.APIKey{
	// 	ID:     "1",
	// 	Name:   "Test App",
	// 	Key:    "dk_test_12345",
	// 	Status: apikey.KeyStatusActive,
	// 	ExpiresAt: time.Now().AddDate(1, 0, 0), // Expires in 1 year
	// 	RateLimit: 1000,
	// }
	// apiKeyService.AddKey(testKey)
	// Create a group for children routes with auth middleware
	childrenGroup := app.Group("/api/children", apiKeyMiddleware.ValidateAPIKey())

	// Routes
	childrenGroup.Post("/", handler.CreateChild)
	childrenGroup.Get("/", handler.GetChildren)
	childrenGroup.Get("/:id", handler.GetChild)
	childrenGroup.Put("/:id", handler.UpdateChild)
}
