package repository

import (
	"dailyalu-server/internal/module/user/domain"
	"time"
)

// UserRepository defines the interface for user data access
type IUserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	UpdateLastLogin(id string, lastLogin time.Time) error
	GetByVerificationToken(token string) (*domain.User, error)
	UpdatePassword(id, password string) error
}
