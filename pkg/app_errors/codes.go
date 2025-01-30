package app_errors

const (
	// Common Errors (4000-4099)
	ErrCodeBadRequest     ErrorCode = 4000
	ErrCodeUnauthorized   ErrorCode = 4001
	ErrCodeForbidden      ErrorCode = 4002
	ErrCodeNotFound       ErrorCode = 4003
	ErrCodeValidation     ErrorCode = 4004
	ErrCodeRateLimit      ErrorCode = 4005
	ErrCodeInvalidInput   ErrorCode = 4006

	// Authentication Errors (4100-4199)
	ErrCodeInvalidCredentials    ErrorCode = 4100
	ErrCodeTokenExpired         ErrorCode = 4101
	ErrCodeInvalidToken         ErrorCode = 4102
	ErrCodeInvalidRefreshToken  ErrorCode = 4103

	// Server Errors (5000-5099)
	ErrCodeInternal      ErrorCode = 5000
	ErrCodeDatabase      ErrorCode = 5001
	ErrCodeThirdParty    ErrorCode = 5002
	ErrCodeConfiguration ErrorCode = 5003
)

// Error messages
var errorMessages = map[ErrorCode]string{
	// 4xxx Client Errors
	ErrCodeBadRequest:          "Bad request",
	ErrCodeUnauthorized:        "Unauthorized",
	ErrCodeForbidden:           "Forbidden",
	ErrCodeNotFound:            "Resource not found",
	ErrCodeValidation:          "Validation failed",
	ErrCodeRateLimit:           "Rate limit exceeded",
	ErrCodeInvalidInput:        "Invalid input",
	
	// Authentication Errors
	ErrCodeInvalidCredentials:   "Invalid credentials",
	ErrCodeTokenExpired:        "Token has expired",
	ErrCodeInvalidToken:        "Invalid token",
	ErrCodeInvalidRefreshToken: "Invalid refresh token",

	// 5xxx Server Errors
	ErrCodeInternal:           "Internal server error",
	ErrCodeDatabase:           "Database error",
	ErrCodeThirdParty:         "Third party service error",
	ErrCodeConfiguration:      "Configuration error",
}

// GetErrorMessage returns the default message for an error code
func GetErrorMessage(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}
