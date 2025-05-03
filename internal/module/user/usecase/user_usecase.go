package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/security/password"
	"dailyalu-server/internal/security/token"
	"dailyalu-server/internal/service/mailer"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type userUseCase struct {
	repo         repository.IUserRepository
	jwtManager   *jwt.JWTManager
	tokenService *token.TokenService
	mailerService *mailer.MailerService
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo repository.IUserRepository, jwtManager *jwt.JWTManager, tokenService *token.TokenService, mailerService *mailer.MailerService) IUserUseCase {
	return &userUseCase{
		repo:         repo,
		jwtManager:   jwtManager,
		tokenService: tokenService,
		mailerService: mailerService,
	}
}

// Implementation of UserUseCase interface
func (uc *userUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	if req.Password != req.ConfirmPassword {
		return nil, ErrDifferentConfirmationPassword
	}
	
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
	verificationToken, err := uc.tokenService.GenerateToken()
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

	verificationLink := uc.tokenService.GenerateVerificationLink("https://dailyalu.mom/verify-email", verificationToken)
	
	// Implement email sending logic here
	_, err = uc.mailerService.Send(ctx, "noreply@dailyalu.mom", req.Email, "Email Verification Link", "Hello, thank you for registering to Daily Alu. Verify your account in Daily Alu by clicking the following url "+verificationLink, "")

	if err != nil {
		return nil, fmt.Errorf("failed to send email")
	}

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

	if uc.tokenService.IsTokenExpired("email", user.CreatedAt) {
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

func (uc *userUseCase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := uc.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	accessToken, refreshToken, err := uc.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Return login response
	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (uc *userUseCase) RefreshToken(refreshToken string) (string, string, error) {
	// Use JWT manager to validate and generate new tokens
	accessToken, newRefreshToken, err := uc.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

func (uc *userUseCase) GetUser(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *userUseCase) UpdateUser(request *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := uc.repo.GetByID(request.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Email = request.Email
	user.Name = request.Name
	user.UpdatedAt = time.Now()

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) UpdatePassword(request *domain.UpdatePasswordRequest) error {
  if request.NewPassword != request.ConfirmPassword {
    return ErrDifferentConfirmationPassword
  }
  
  user, err := uc.repo.GetByID(request.ID)
  if err != nil {
    return err
  }

  if user == nil {
    return ErrUserNotFound
  }

  if !password.Verify(request.OldPassword, user.PasswordHash) {
		return ErrInvalidOldPassword
	}

  hashedPassword, err := password.Hash(request.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

  err = uc.repo.UpdatePassword(request.ID, hashedPassword)
  if err != nil {
    return err
  }

  return nil
}

func (uc *userUseCase) DeleteUser(id string) error {
	return uc.repo.Delete(id)
}

func (uc *userUseCase) ForgotPassword(req *domain.ForgotPasswordRequest) error {
    // Find user by email
    user, err := uc.repo.GetByEmail(req.Email)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    if user == nil {
        // For security, don't reveal if email exists
        return nil
    }

    // Generate reset token
    resetToken, err := uc.tokenService.GenerateToken()
    if err != nil {
        return fmt.Errorf("failed to generate reset token: %w", err)
    }

    // Update user with reset token
    if err := uc.repo.UpdateForgotPasswordToken(user.ID, resetToken); err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    // Generate reset link
    //resetLink := uc.tokenService.GeneratePasswordResetLink("http://your-frontend-url", resetToken)
    
    // TODO: Send password reset email with resetLink
    // Implement email sending logic here

    return nil
}

func (uc *userUseCase) ResetPassword(req *domain.ResetPasswordRequest) error {
    if req.NewPassword != req.ConfirmPassword {
        return ErrDifferentConfirmationPassword
    }

    // Find user by reset token
    user, err := uc.repo.GetByResetPasswordToken(req.Token)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    if user == nil {
        return ErrInvalidResetToken
    }

    // Check token expiration
    if uc.tokenService.IsTokenExpired(token.PasswordReset, user.ResetPasswordTokenRequestedAt) {
        return ErrResetTokenExpired
    }

    // Hash new password
    hashedPassword, err := password.Hash(req.NewPassword)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Update user password and clear reset token
    user.PasswordHash = hashedPassword
    user.ResetPasswordToken = ""
    user.UpdatedAt = time.Now()
    
    if err := uc.repo.Update(user); err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }

    return nil
}