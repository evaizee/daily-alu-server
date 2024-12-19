package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
)

// UserUseCase defines the interface for user business logic
type IUserUseCase interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
	Login(req *domain.LoginRequest) (string, string, error)
	GetUser(id string) (*domain.User, error)
	UpdateUser(id, email, name string) (*domain.User, error)
	DeleteUser(id string) error
	VerifyEmail(ctx context.Context, token string) error
	RefreshToken(refreshToken string) (string, string, error)
}