package usecase

import "errors"

// Domain errors for children module
var (
	ErrChildNotFound      = errors.New("child not found")
	ErrUnauthorizedAccess = errors.New("unauthorized access to child data")
	ErrInvalidChildData   = errors.New("invalid child data")
)
