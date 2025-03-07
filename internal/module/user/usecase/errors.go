package usecase

import "errors"

var (
	ErrUserNotFound                  = errors.New("User not found")
	ErrUserAlreadyExists             = errors.New("User already exists")
	ErrEmailAlreadyExists            = errors.New("Email already exists")
	ErrInvalidCredentials            = errors.New("Invalid credentials")
	ErrInvalidVerificationToken      = errors.New("Invalid verification token")
	ErrVerificationTokenExpired      = errors.New("Token expired")
	ErrDifferentConfirmationPassword = errors.New("Confirmation password is not equal")
	ErrInvalidOldPassword            = errors.New("Old password is not valid")
)
