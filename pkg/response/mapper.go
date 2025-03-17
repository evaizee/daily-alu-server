package response

import (
	"dailyalu-server/internal/module/user/usecase"
	"errors"
)

// MapDomainError maps domain-specific errors to standardized API errors
func MapDomainError(err error) *AppError {
	// User domain errors
	switch {
	case errors.Is(err, usecase.ErrEmailAlreadyExists):
		return NewBadRequestError("Email already exists")
	case errors.Is(err, usecase.ErrInvalidCredentials):
		return NewUnauthorizedError("Invalid credentials")
	case errors.Is(err, usecase.ErrUserNotFound):
		return NewNotFoundError("User not found")
	case errors.Is(err, usecase.ErrDifferentConfirmationPassword):
		return NewValidationError("Password and confirmation password do not match")
	case errors.Is(err, usecase.ErrInvalidVerificationToken):
		return NewBadRequestError("Invalid verification token")
	case errors.Is(err, usecase.ErrVerificationTokenExpired):
		return NewBadRequestError("Verification token has expired")
	case errors.Is(err, usecase.ErrInvalidResetToken):
		return NewBadRequestError("Invalid password reset token")
	case errors.Is(err, usecase.ErrResetTokenExpired):
		return NewBadRequestError("Password reset token has expired")
	case errors.Is(err, usecase.ErrInvalidOldPassword):
		return NewBadRequestError("Invalid old password")
	// Add mappings for other domain errors as needed
	
	// Default case - internal error
	default:
		return NewInternalError(err)
	}
}
