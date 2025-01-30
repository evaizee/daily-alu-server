package usecase

import "errors"

var (
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserAlreadyExists             = errors.New("user already exists")
	ErrEmailAlreadyExists            = errors.New("email already exists")
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrInvalidVerificationToken      = errors.New("invalid verification token")
	ErrVerificationTokenExpired      = errors.New("token expired")
	ErrDifferentConfirmationPassword = errors.New("confirmation password is not equal")
)
