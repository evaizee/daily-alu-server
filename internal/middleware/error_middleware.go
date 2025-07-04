package middleware

import (
	"dailyalu-server/pkg/app_log/zap_log"
	"dailyalu-server/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorMiddleware struct {
	logger *zap.Logger
}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{
		logger: zap_log.Logger,
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
		
		var bodyCopy string
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPatch {
			bodyBytes := c.Body()
			if len(bodyBytes) > 0 && len(bodyBytes) < 5*1024 { // limit to 5KB
				bodyCopy = string(bodyBytes)
			} else if len(bodyBytes) >= 5*1024 {
				bodyCopy = "[body too large to log]"
			}
		}
		
		// Convert to our AppError type if it isn't already
		var appErr *response.AppError
		if appError, ok := err.(*response.AppError); ok {
			appErr = appError
		} else {
			// Unknown error, convert to internal error
			appErr = response.NewInternalError(err)
		}

		// Prepare log fields
		fields := []zap.Field{
			zap.String("request_id", reqID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", c.IP()),
			zap.String("body", bodyCopy),
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

		// Return error response
		return response.Error(c, appErr)
	}
}
