package middleware

import (
	"dailyalu-server/pkg/response"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/spf13/viper"
)

// EndpointID generates a consistent identifier for an endpoint
func EndpointID(method, path string) string {
	// Normalize method to lowercase
	method = strings.ToLower(method)
	
	// Convert path to identifier by replacing slashes with underscores
	// and removing any path parameters
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, ":", "")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	
	return method + "." + path
}

// CreateEndpointRateLimiter creates a rate limiter specific to an endpoint
func CreateEndpointRateLimiter(method, path string) fiber.Handler {
	// Skip if rate limiting is disabled globally
	if !viper.GetBool("ratelimit.enabled") {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	endpointID := EndpointID(method, path)
	configPath := fmt.Sprintf("ratelimit.endpoints.%s", endpointID)
	
	// Check if endpoint-specific config exists
	var max int
	var expiration time.Duration
	
	if viper.IsSet(configPath + ".max") {
		// Use endpoint-specific settings
		max = viper.GetInt(configPath + ".max")
		expiration = time.Duration(viper.GetInt(configPath + ".expiration")) * time.Second
	} else {
		// Fall back to default settings
		max = viper.GetInt("ratelimit.default.max")
		expiration = time.Duration(viper.GetInt("ratelimit.default.expiration")) * time.Second
	}
	
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		// Use the endpoint ID in the key to ensure separate limits per endpoint
		KeyGenerator: func(c *fiber.Ctx) string {
			return endpointID + ":" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, response.NewTooManyRequestsError(
				fmt.Sprintf("Rate limit exceeded for this endpoint. Maximum %d requests per %d seconds.", 
					max, int(expiration.Seconds()))))
		},
	})
}

// RateLimitedRoute applies a rate limiter to a route based on method and path
func RateLimitedRoute(app fiber.Router, method, path string, handlers ...fiber.Handler) fiber.Router {
	 // Get the router group's prefix if available
	 var fullPath string
    
	 // Check if the app is a group with a prefix
	 if group, ok := app.(*fiber.Group); ok {
		 // Combine the group prefix with the path
		 fullPath = group.Prefix + path
	 } else {
		 // If not a group, just use the path as is
		 fullPath = path
	 }

	// Create the rate limiter for this endpoint
	rateLimiter := CreateEndpointRateLimiter(method, fullPath)
	
	// Prepend the rate limiter to the handlers
	allHandlers := append([]fiber.Handler{rateLimiter}, handlers...)
	
	// Register the route with the appropriate method
	switch strings.ToUpper(method) {
	case "GET":
		return app.Get(path, allHandlers...)
	case "POST":
		return app.Post(path, allHandlers...)
	case "PUT":
		return app.Put(path, allHandlers...)
	case "DELETE":
		return app.Delete(path, allHandlers...)
	case "PATCH":
		return app.Patch(path, allHandlers...)
	default:
		return app.All(path, allHandlers...)
	}
}
