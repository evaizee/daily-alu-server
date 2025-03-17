package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
)

// UserUseCase defines the interface for user business logic
type IUserUseCase interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
	Login(req *domain.LoginRequest) (*domain.LoginResponse, error)
	GetUser(id string) (*domain.User, error)
	UpdateUser(req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(id string) error
	VerifyEmail(ctx context.Context, token string) error
	RefreshToken(refreshToken string) (string, string, error)
	UpdatePassword(request *domain.UpdatePasswordRequest) error
	ForgotPassword(req *domain.ForgotPasswordRequest) error
	ResetPassword(req *domain.ResetPasswordRequest) error
}