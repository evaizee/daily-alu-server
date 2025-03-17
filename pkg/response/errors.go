package response

// NewAppError creates a new AppError
func NewAppError(errType ErrorType, code int, message string) *AppError {
	if message == "" {
		message = GetErrorMessage(code)
	}
	
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrorTypeClient, ErrCodeBadRequest, message)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrorTypeClient, ErrCodeUnauthorized, message)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrorTypeClient, ErrCodeForbidden, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrorTypeClient, ErrCodeNotFound, message)
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return NewAppError(ErrorTypeClient, ErrCodeValidation, message)
}

// NewValidationErrorWithDetails creates a validation error with details
func NewValidationErrorWithDetails(message string, details interface{}) *AppError {
	err := NewValidationError(message)
	err.Details = details
	return err
}

// NewInternalError creates an internal server error
func NewInternalError(err error) *AppError {
	appErr := NewAppError(ErrorTypeServer, ErrCodeInternal, "")
	if err != nil {
		appErr.Internal = err
	}
	return appErr
}

// NewDatabaseError creates a database error
func NewDatabaseError(err error) *AppError {
	appErr := NewAppError(ErrorTypeServer, ErrCodeDatabase, "")
	if err != nil {
		appErr.Internal = err
	}
	return appErr
}
