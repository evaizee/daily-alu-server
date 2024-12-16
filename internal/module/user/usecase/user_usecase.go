package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/security/password"
	"dailyalu-server/internal/service/email"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type userUseCase struct {
	repo               repository.IUserRepository
	jwtManager         *jwt.JWTManager
	verificationService *email.VerificationService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo repository.IUserRepository, jwtManager *jwt.JWTManager, verificationService *email.VerificationService) IUserUseCase {
	return &userUseCase{
		repo:               repo,
		jwtManager:         jwtManager,
		verificationService: verificationService,
	}
}

// Implementation of UserUseCase interface
func (uc *userUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := uc.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	verificationToken, err := uc.verificationService.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:                   uuid.New().String(),
		Email:                req.Email,
		Name:                 req.Name,
		PasswordHash:         hashedPassword,
		Status:               domain.UserStatusNotActive,
		EmailVerificationToken: verificationToken,
		Role:                 "user",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Send verification email
	//verificationLink := uc.verificationService.GenerateVerificationLink("http://your-api-domain", verificationToken)
	// Implement email sending logic here

	return user, nil
}

func (uc *userUseCase) VerifyEmail(ctx context.Context, token string) error {
	user, err := uc.repo.GetByVerificationToken(token)
	if err != nil {
		return fmt.Errorf("failed to get user by token: %w", err)
	}
	if user == nil {
		return ErrInvalidVerificationToken
	}

	if uc.verificationService.IsTokenExpired(user.CreatedAt) {
		return ErrVerificationTokenExpired
	}

	user.Status = domain.UserStatusActive
	user.EmailVerificationToken = "" // Clear the token after verification
	user.UpdatedAt = time.Now()

	if err := uc.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	return nil
}

func (uc *userUseCase) Login(req *domain.LoginRequest) (string, error) {
	user, err := uc.repo.GetByEmail(req.Email)
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
	if err := uc.repo.UpdateLastLogin(user.ID, now); err != nil {
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
	return uc.repo.GetByID(id)
}

func (uc *userUseCase) UpdateUser(id, email, name string) (*domain.User, error) {
	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Email = email
	user.Name = name
	user.UpdatedAt = time.Now()

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}

// Helper function to generate unique ID
func generateID() string {
	return time.Now().Format("20060102150405") // Simple ID generation for demo
}
