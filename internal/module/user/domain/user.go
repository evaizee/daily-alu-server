package domain

import (
	"time"
)

// User status constants
const (
	UserStatusNotActive = 0
	UserStatusActive    = 10
	UserStatusBlocked   = 20
)

type User struct {
	ID                            string     `json:"id"`
	Email                         string     `json:"email"`
	Name                          string     `json:"name"`
	PasswordHash                  string     `json:"-"`
	Status                        int16      `json:"-"`
	EmailVerificationToken        string     `json:"-"`
	ResetPasswordToken            string     `json:"-"`
	ResetPasswordTokenRequestedAt time.Time  `json:"-"`
	Role                          string     `json:"-"`
	LastLogin                     *time.Time `json:"last_login,omitempty"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsBlocked checks if the user is blocked
func (u *User) IsBlocked() bool {
	return u.Status == UserStatusBlocked
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	Name            string `json:"name" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdatePasswordRequest struct {
	ID              string `json:"id"`
	OldPassword     string `json:"old_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	ID    string `json:"id"`
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
