package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/security/token"
	"testing"
	"time"
)

func TestUserUseCase_VerifyEmail(t *testing.T) {
	testCases := []struct {
		name string
		token          string
		mockRepo      *MockUserRepository
		mockToken     *token.TokenService
		mockMailer    *MockMailerService
		expectedError error
	}{
		{
			name: "verify email success",
			token: "abcdefghijk",
			mockToken: token.NewTokenService(),
			mockRepo: &MockUserRepository{
				GetByVerificationTokenFunc: func(token string) (*domain.User, error) {
					return &domain.User{
						Email:                  "test@example.com",
						Name:                   "Test User",
						Status:                 domain.UserStatusNotActive,
						EmailVerificationToken: "verification-token",
						Role:                   "user",
						CreatedAt: time.Now(),
					}, nil
				},
				UpdateFunc: func(user *domain.User) error {
					return nil
				},
			},
			expectedError: nil,
		},
		{
			name: "verification token not found",
			token: "abcdefghijk",
			mockToken: token.NewTokenService(),
			mockRepo: &MockUserRepository{
				GetByVerificationTokenFunc: func(token string) (*domain.User, error) {
					return nil, nil
				},
				UpdateFunc: func(user *domain.User) error {
					return nil
				},
			},
			expectedError: ErrInvalidVerificationToken,
		},
		{
			name: "verification token expired",
			token: "abcdefghijk",
			mockToken: token.NewTokenService(),
			mockRepo: &MockUserRepository{
				GetByVerificationTokenFunc: func(token string) (*domain.User, error) {
					return &domain.User{
						Email:                  "test@example.com",
						Name:                   "Test User",
						Status:                 domain.UserStatusNotActive,
						EmailVerificationToken: "verification-token",
						Role:                   "user",
						CreatedAt: time.Date(2001, 10, 1, 12, 0, 0, 0, time.Local),
					}, nil
				},
				UpdateFunc: func(user *domain.User) error {
					return nil
				},
			},
			expectedError: ErrVerificationTokenExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := &userUseCase{
				repo:         tc.mockRepo,
				tokenService: tc.mockToken,
				mailerService: tc.mockMailer,
			}

			err := uc.VerifyEmail(context.Background(), tc.token)

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
		})
	}
}