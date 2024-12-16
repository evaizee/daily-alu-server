package container

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/service/email"
	"database/sql"
	"time"

	//"github.com/docker/docker/volume/service"
)

// Container holds all the dependencies for the application
type Container struct {
	// Infrastructure
	db         *sql.DB
	jwtManager *jwt.JWTManager

	// Repositories
	userRepository repository.IUserRepository

	// Use Cases
	userUseCase usecase.IUserUseCase

	// Handlers
	userHandler *api.UserHandler

	// Middleware
	securityMiddleware *middleware.SecurityMiddleware

	//Email Verification Service
	emailVerificationService *email.VerificationService
}

// NewContainer creates a new dependency injection container
func NewContainer(db *sql.DB, jwtSecret string, jwtExpiry time.Duration) *Container {
	c := &Container{
		db: db,
	}

	// Initialize JWT manager
	c.jwtManager = jwt.NewJWTManager(jwtSecret, jwtExpiry)

	// Initialize repositories
	c.userRepository = repository.NewPostgresUserRepository(db)

	c.emailVerificationService = email.NewVerificationService()

	// Initialize use cases
	c.userUseCase = usecase.NewUserUseCase(c.userRepository, c.jwtManager, c.emailVerificationService)

	// Initialize handlers
	c.userHandler = api.NewUserHandler(c.userUseCase)

	// Initialize middleware
	c.securityMiddleware = middleware.NewSecurityMiddleware(middleware.SecurityConfig{
		JWTManager: c.jwtManager,
	})

	return c
}

// GetUserHandler returns the user handler
func (c *Container) GetUserHandler() *api.UserHandler {
	return c.userHandler
}

// GetSecurityMiddleware returns the security middleware
func (c *Container) GetSecurityMiddleware() *middleware.SecurityMiddleware {
	return c.securityMiddleware
}

// Close closes any resources held by the container
func (c *Container) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
