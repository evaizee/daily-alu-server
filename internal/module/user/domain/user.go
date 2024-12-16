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
	ID                    string     `json:"id"`
	Email                 string     `json:"email"`
	Name                  string     `json:"name"`
	PasswordHash          string     `json:"-"`
	Status                int16      `json:"status"`
	EmailVerificationToken string     `json:"-"`
	Role                  string     `json:"role"`
	LastLogin             *time.Time `json:"last_login,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}
