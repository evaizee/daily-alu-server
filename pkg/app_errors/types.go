package app_errors

import "fmt"

// ErrorCode represents a unique error code
type ErrorCode int

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeClient ErrorType = "client" // 4xx errors
	ErrorTypeServer ErrorType = "server" // 5xx errors
)

// AppError represents an application error
type AppError struct {
	Type     ErrorType              `json:"-"`
	Code     ErrorCode             `json:"code"`
	Message  string                `json:"message"`
	Internal error                 `json:"-"`
	Metadata map[string]interface{} `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Response formats the error for API response
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents the error structure sent to clients
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Response converts AppError to ErrorResponse
func (e *AppError) Response() ErrorResponse {
	return ErrorResponse{
		Error: APIError{
			Code:    e.Code,
			Message: e.Message,
		},
	}
}

// WithInternal adds internal error details
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(metadata map[string]interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// AddMetadata adds a single metadata key-value pair
func (e *AppError) AddMetadata(key string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}
