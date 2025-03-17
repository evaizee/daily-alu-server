package response

import (
	"github.com/gofiber/fiber/v2"
)

// Success sends a standardized success response
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(StandardResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error response
func Error(c *fiber.Ctx, appErr *AppError) error {
	status := getHTTPStatus(appErr.Code)
	return c.Status(status).JSON(appErr.Response())
}

// getHTTPStatus maps error codes to HTTP status codes
func getHTTPStatus(code int) int {
	switch {
	case code >= 5000:
		return fiber.StatusInternalServerError
	case code >= 4100:
		return fiber.StatusUnauthorized
	case code >= 4050:
		return fiber.StatusMethodNotAllowed
	case code >= 4040:
		return fiber.StatusNotFound
	case code >= 4030:
		return fiber.StatusForbidden
	case code >= 4020:
		return fiber.StatusUnprocessableEntity
	case code >= 4010:
		return fiber.StatusUnauthorized
	case code >= 4000:
		return fiber.StatusBadRequest
	default:
		return fiber.StatusInternalServerError
	}
}
