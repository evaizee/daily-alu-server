package usecase

import (
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/security/password"
	"time"
)

type userUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) IUserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Implementation of UserUseCase interface
func (uc *userUseCase) Register(req domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           generateID(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		Role:         "user", // Default role
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Login(req domain.LoginRequest) (string, error) {
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if !password.Verify(req.Password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	if err := uc.userRepo.UpdateLastLogin(user.ID, now); err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := uc.jwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *userUseCase) UpdateUser(id, email, name string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Email = email
	user.Name = name
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	return uc.userRepo.Delete(id)
}

// Helper function to generate unique ID
func generateID() string {
	return time.Now().Format("20060102150405") // Simple ID generation for demo
}
