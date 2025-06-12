package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/security/token"
	"errors"
	"testing"
)

func TestUserUseCase_Register(t *testing.T) {
	testCases := []struct {
		name          string
		req           *domain.RegisterRequest
		mockRepo      *MockUserRepository
		mockToken     *token.TokenService
		mockMailer    *MockMailerService
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name: "successful registration",
			req: &domain.RegisterRequest{
				Email:           "test@example.com",
				Name:            "Test User",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			mockRepo: &MockUserRepository{
				GetByEmailFunc: func(email string) (*domain.User, error) {
					return nil, nil // user doesn't exist
				},
				CreateFunc: func(user *domain.User) error {
					return nil // successful creation
				},
			},
			mockMailer: &MockMailerService{
				SendVerificationEmailFunc: func() error {
					return nil
				},
			},
			expectedUser: &domain.User{
				Email:                  "test@example.com",
				Name:                   "Test User",
				Status:                 domain.UserStatusNotActive,
				EmailVerificationToken: "verification-token",
				Role:                   "user",
			},
			expectedError: nil,
		},
		{
			name: "password mismatch",
			req: &domain.RegisterRequest{
				Email:           "test@example.com",
				Name:            "Test User",
				Password:        "password123",
				ConfirmPassword: "different",
			},
			expectedError: ErrDifferentConfirmationPassword,
		},
		{
			name: "email already exists",
			req: &domain.RegisterRequest{
				Email:           "existing@example.com",
				Name:            "Test User",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			mockRepo: &MockUserRepository{
				GetByEmailFunc: func(email string) (*domain.User, error) {
					return &domain.User{Email: email}, nil // user exists
				},
			},
			expectedError: ErrEmailAlreadyExists,
		},
		{
			name: "repository error when checking email",
			req: &domain.RegisterRequest{
				Email:           "test@example.com",
				Name:            "Test User",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			mockRepo: &MockUserRepository{
				GetByEmailFunc: func(email string) (*domain.User, error) {
					return nil, errors.New("database error")
				},
			},
			expectedError: errors.New("failed to check existing user: database error"),
		},
		// {
		// 	name: "password hashing error",
		// 	req: &domain.RegisterRequest{
		// 		Email:           "test@example.com",
		// 		Name:            "Test User",
		// 		Password:        "password123",
		// 		ConfirmPassword: "password123",
		// 	},
		// 	mockRepo: &MockUserRepository{
		// 		GetByEmailFunc: func(email string) (*domain.User, error) {
		// 			return nil, nil
		// 		},
		// 	},
		// 	expectedError: errors.New("failed to hash password"), // This will need adjustment
		// },
		// {
		// 	name: "token generation error",
		// 	req: &domain.RegisterRequest{
		// 		Email:           "test@example.com",
		// 		Name:            "Test User",
		// 		Password:        "password123",
		// 		ConfirmPassword: "password123",
		// 	},
		// 	mockRepo: &MockUserRepository{
		// 		GetByEmailFunc: func(email string) (*domain.User, error) {
		// 			return nil, nil
		// 		},
		// 	},
		// 	expectedError: errors.New("failed to generate verification token: token generation failed"),
		// },
		{
			name: "user creation error",
			req: &domain.RegisterRequest{
				Email:           "test@example.com",
				Name:            "Test User",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			mockRepo: &MockUserRepository{
				GetByEmailFunc: func(email string) (*domain.User, error) {
					return nil, nil
				},
				CreateFunc: func(user *domain.User) error {
					return errors.New("creation failed")
				},
			},
			expectedError: errors.New("failed to create user: creation failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := &userUseCase{
				repo:         tc.mockRepo,
				tokenService: tc.mockToken,
				mailerService: tc.mockMailer,
			}

			user, err := uc.Register(context.Background(), tc.req)

			// Error assertion
			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tc.expectedError)
					return
				}
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// User assertion
			if user == nil {
				t.Error("expected user, got nil")
				return
			}

			if user.Email != tc.req.Email {
				t.Errorf("expected email %s, got %s", tc.req.Email, user.Email)
			}

			if user.Name != tc.req.Name {
				t.Errorf("expected name %s, got %s", tc.req.Name, user.Name)
			}

			if user.Status != domain.UserStatusNotActive {
				t.Errorf("status not match")
			}

			if user.Role != "user" {
				t.Errorf("expected role 'user', got %s", user.Role)
			}

			if user.EmailVerificationToken == "" {
				t.Error("expected verification token, got empty")
			}

			if user.PasswordHash == "" {
				t.Error("expected password hash, got empty")
			}
		})
	}
}