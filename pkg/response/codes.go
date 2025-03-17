package response

// Error codes
const (
	// Common Errors (4000-4099)
	ErrCodeBadRequest   = 4000
	ErrCodeUnauthorized = 4001
	ErrCodeForbidden    = 4002
	ErrCodeNotFound     = 4003
	ErrCodeValidation   = 4004
	ErrCodeRateLimit    = 4005
	ErrCodeInvalidInput = 4006

	// Authentication Errors (4100-4199)
	ErrCodeInvalidCredentials  = 4100
	ErrCodeTokenExpired        = 4101
	ErrCodeInvalidToken        = 4102
	ErrCodeInvalidRefreshToken = 4103

	// Server Errors (5000-5099)
	ErrCodeInternal      = 5000
	ErrCodeDatabase      = 5001
	ErrCodeThirdParty    = 5002
	ErrCodeConfiguration = 5003
)

// Error messages
var errorMessages = map[int]string{
	// 4xxx Client Errors
	ErrCodeBadRequest:   "Bad request",
	ErrCodeUnauthorized: "Unauthorized",
	ErrCodeForbidden:    "Forbidden",
	ErrCodeNotFound:     "Resource not found",
	ErrCodeValidation:   "Validation failed",
	ErrCodeRateLimit:    "Rate limit exceeded",
	ErrCodeInvalidInput: "Invalid input",

	// Authentication Errors
	ErrCodeInvalidCredentials:  "Invalid credentials",
	ErrCodeTokenExpired:        "Token has expired",
	ErrCodeInvalidToken:        "Invalid token",
	ErrCodeInvalidRefreshToken: "Invalid refresh token",

	// 5xxx Server Errors
	ErrCodeInternal:      "Internal server error",
	ErrCodeDatabase:      "Database error",
	ErrCodeThirdParty:    "Third party service error",
	ErrCodeConfiguration: "Configuration error",
}

// GetErrorMessage returns the default message for an error code
func GetErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}
