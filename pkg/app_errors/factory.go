package app_errors

// NewError creates a new AppError
func NewError(code ErrorCode, message string) *AppError {
	errType := ErrorTypeServer
	if code < 5000 {
		errType = ErrorTypeClient
	}

	if message == "" {
		message = GetErrorMessage(code)
	}

	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

// Common error constructors
func NewBadRequestError(message string) *AppError {
	return NewError(ErrCodeBadRequest, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewError(ErrCodeUnauthorized, message)
}

func NewForbiddenError(message string) *AppError {
	return NewError(ErrCodeForbidden, message)
}

func NewNotFoundError(message string) *AppError {
	return NewError(ErrCodeNotFound, message)
}

func NewValidationError(message string) *AppError {
	return NewError(ErrCodeValidation, message)
}

func NewInternalError(err error) *AppError {
	return NewError(ErrCodeInternal, "Internal server error").
		WithInternal(err)
}

func NewDatabaseError(err error) *AppError {
	return NewError(ErrCodeDatabase, "Database error").
		WithInternal(err)
}

// Authentication error constructors
func NewInvalidCredentialsError() *AppError {
	return NewError(ErrCodeInvalidCredentials, "")
}

func NewTokenExpiredError() *AppError {
	return NewError(ErrCodeTokenExpired, "")
}

func NewInvalidTokenError() *AppError {
	return NewError(ErrCodeInvalidToken, "")
}
