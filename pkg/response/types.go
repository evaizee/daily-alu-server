package response

import "fmt"

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeClient ErrorType = "client" // 4xx errors
	ErrorTypeServer ErrorType = "server" // 5xx errors
)

// StandardResponse is the standard API response structure for success cases
type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard API response structure for error cases
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
	NextPage     *int  `json:"next_page,omitempty"`
	PreviousPage *int  `json:"previous_page,omitempty"`
}

// PaginatedResponse is the standard API response structure with pagination
type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// AppError represents an application error
type AppError struct {
	Type     ErrorType              `json:"-"`
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Internal error                  `json:"-"`
	Metadata map[string]interface{} `json:"-"`
	Details  interface{}            `json:"details,omitempty"` // For validation errors
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Response formats the error for API response
func (e *AppError) Response() interface{} {
	if e.Details != nil {
		// Return validation error format with details
		return ErrorResponse{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	}
	
	// Return standard error format
	return ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}

// WithDetails adds validation details to the error
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
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
