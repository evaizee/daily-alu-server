package response

import (
	childrenUsecase "dailyalu-server/internal/module/children/usecase"
	userUsecase "dailyalu-server/internal/module/user/usecase"
	"errors"
)

// MapDomainError maps domain-specific errors to standardized API errors
func MapDomainError(err error) *AppError {
	// User domain errors
	switch {
	case errors.Is(err, userUsecase.ErrEmailAlreadyExists):
		return NewBadRequestError("Email already exists")
	case errors.Is(err, userUsecase.ErrInvalidCredentials):
		return NewUnauthorizedError("Invalid credentials")
	case errors.Is(err, userUsecase.ErrUserNotFound):
		return NewNotFoundError("User not found")
	case errors.Is(err, userUsecase.ErrDifferentConfirmationPassword):
		return NewValidationError("Password and confirmation password do not match")
	case errors.Is(err, userUsecase.ErrInvalidVerificationToken):
		return NewBadRequestError("Invalid verification token")
	case errors.Is(err, userUsecase.ErrVerificationTokenExpired):
		return NewBadRequestError("Verification token has expired")
	case errors.Is(err, userUsecase.ErrInvalidResetToken):
		return NewBadRequestError("Invalid password reset token")
	case errors.Is(err, userUsecase.ErrResetTokenExpired):
		return NewBadRequestError("Password reset token has expired")
	case errors.Is(err, userUsecase.ErrInvalidOldPassword):
		return NewBadRequestError("Invalid old password")
	
	// Children domain errors
	case errors.Is(err, childrenUsecase.ErrChildNotFound):
		return NewNotFoundError("Child not found")
	case errors.Is(err, childrenUsecase.ErrUnauthorizedAccess):
		return NewForbiddenError("You do not have permission to access this child's data")
	case errors.Is(err, childrenUsecase.ErrInvalidChildData):
		return NewBadRequestError("Invalid child data")
	
	// Default case - internal error
	default:
		return NewInternalError(err)
	}
}
