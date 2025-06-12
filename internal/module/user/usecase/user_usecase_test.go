// user_usecase_test.go
package usecase

import (
	"context"
	"dailyalu-server/internal/module/user/domain"
	mailerDomain "dailyalu-server/internal/service/mailer/domain"
	"time"
)

// MockUserRepository implements the repository interface for testing
type MockUserRepository struct {
	GetByEmailFunc                func(email string) (*domain.User, error)
	GetByIDFunc func(id string) (*domain.User, error)
	CreateFunc                    func(user *domain.User) error
	UpdateFunc                    func(user *domain.User) error
	DeleteFunc                    func(id string) error
	UpdateLastLoginFunc           func(id string, lastLogin time.Time) error
	GetByVerificationTokenFunc    func(token string) (*domain.User, error)
	UpdatePasswordFunc            func(id, password string) error
	GetByResetPasswordTokenFunc   func(token string) (*domain.User, error)
	UpdateForgotPasswordTokenFunc func(id, token string) error
}

func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	return m.GetByIDFunc(id)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	return m.GetByEmailFunc(email)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	return m.CreateFunc(user)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	return m.UpdateFunc(user)
}

func (m *MockUserRepository) Delete(id string) error {
	return m.DeleteFunc(id)
}

func (m *MockUserRepository) UpdateLastLogin(id string, lastLogin time.Time) error {
	return m.UpdateLastLoginFunc(id, lastLogin)
}

func (m *MockUserRepository) GetByVerificationToken(token string) (*domain.User, error) {
	return m.GetByVerificationTokenFunc(token)
}

func (m *MockUserRepository) UpdatePassword(id, password string) error {
	return m.UpdatePasswordFunc(id, password)
}

func (m *MockUserRepository) GetByResetPasswordToken(token string) (*domain.User, error) {
	return m.GetByResetPasswordTokenFunc(token)
}

func (m *MockUserRepository) UpdateForgotPasswordToken(id, token string) error {
	return m.UpdateForgotPasswordTokenFunc(id, token)
}

// MockTokenService implements the token service interface for testing
type MockTokenService struct {
	GenerateTokenFunc            func() (string, error)
	GenerateVerificationLinkFunc func(baseURL, token string) string
}

func (m *MockTokenService) GenerateToken() (string, error) {
	return m.GenerateTokenFunc()
}

func (m *MockTokenService) GenerateVerificationLink(baseURL, token string) string {
	return m.GenerateVerificationLinkFunc(baseURL, token)
}

// MockMailerService implements the mailer service interface for testing
type MockMailerService struct {
	SendVerificationEmailFunc func() error
}

func (m *MockMailerService) SendVerificationEmail(ctx context.Context, data *mailerDomain.EmailVerificationData) error {
	return m.SendVerificationEmailFunc()
}
