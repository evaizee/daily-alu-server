package usecase

import "errors"

var (
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrUserNotFound                  = errors.New("user not found")
	ErrEmailAlreadyExists            = errors.New("email already exists")
	ErrDifferentConfirmationPassword = errors.New("password and confirmation password do not match")
	ErrInvalidVerificationToken      = errors.New("invalid verification token")
	ErrVerificationTokenExpired      = errors.New("verification token has expired")
	ErrInvalidResetToken             = errors.New("invalid password reset token")
	ErrResetTokenExpired             = errors.New("password reset token has expired")
	ErrInvalidOldPassword            = errors.New("invalid old password")
)
