package middleware

import (
	"dailyalu-server/pkg/app_errors"
	"dailyalu-server/pkg/app_log"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorMiddleware struct {
	logger *zap.Logger
}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{
		logger: app_log.Logger,
	}
}

// Handle returns a Fiber middleware function for error handling
func (m *ErrorMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue stack
		err := c.Next()
		if err == nil {
			return nil
		}

		// Get request context and details for logging
		reqID := c.Get("X-Request-ID", "unknown")
		method := c.Method()
		path := c.Path()
		fmt.Println(err)
		// Convert to our AppError type if it isn't already
		var appErr *app_errors.AppError
		if appError, ok := err.(*app_errors.AppError); ok {
			appErr = appError
		} else {
			// Unknown error, convert to internal error
			appErr = app_errors.NewInternalError(err)
		}

		// Prepare log fields
		fields := []zap.Field{
			zap.String("request_id", reqID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", c.IP()),
		}

		// Add error metadata to log fields
		if appErr.Metadata != nil {
			for k, v := range appErr.Metadata {
				fields = append(fields, zap.Any(k, v))
			}
		}

		// Log the error
		if appErr.Internal != nil {
			m.logger.Error("HTTP Error",
				append(fields,
					zap.String("error", appErr.Message),
					zap.Error(appErr.Internal),
				)...,
			)
		} else {
			m.logger.Error("HTTP Error",
				append(fields,
					zap.String("error", appErr.Message),
				)...,
			)
		}

		// Get HTTP status based on error code
		status := getHTTPStatus(appErr.Code)

		// Return error response
		return c.Status(status).JSON(appErr.Response())
	}
}

// getHTTPStatus maps error codes to HTTP status codes
func getHTTPStatus(code app_errors.ErrorCode) int {
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
