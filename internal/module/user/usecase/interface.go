package usecase

import "dailyalu-server/internal/module/user/domain"

// UserUseCase defines the interface for user business logic
type IUserUseCase interface {
	Register(req domain.RegisterRequest) (*domain.User, error)
	Login(req domain.LoginRequest) (string, error)
	GetUser(id string) (*domain.User, error)
	UpdateUser(id, email, name string) (*domain.User, error)
	DeleteUser(id string) error
}
