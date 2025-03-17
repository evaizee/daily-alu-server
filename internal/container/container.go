package container

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	activityRepo "dailyalu-server/internal/module/activity/repository"
	activityUseCase "dailyalu-server/internal/module/activity/usecase"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/security/token"
	"database/sql"
	"time"
)

// Container holds all the dependencies for the application
type Container struct {
	// Infrastructure
	db *sql.DB

	// Managers
	jwtManager *jwt.JWTManager

	// Repositories
	userRepository     repository.IUserRepository
	activityRepository activityRepo.IActivityRepository

	// Use Cases
	userUseCase     usecase.IUserUseCase
	activityUseCase activityUseCase.IActivityUseCase

	// Handlers
	userHandler     *api.UserHandler
	activityHandler *api.ActivityHandler

	// Middleware
	securityMiddleware *middleware.SecurityMiddleware
	errorMiddleware    *middleware.ErrorMiddleware

	//Token Verification Service
	tokenService *token.TokenService
}

// NewContainer creates a new dependency injection container
func NewContainer(db *sql.DB, jwtSecret, jwtRefreshSecretKey string, jwtExpiry, jwtRefreshExpiry time.Duration) *Container {
	c := &Container{
		db: db,
	}

	// Initialize JWT manager
	c.jwtManager = jwt.NewJWTManager(jwtSecret, jwtRefreshSecretKey, jwtExpiry, jwtRefreshExpiry)

	// Initialize repositories
	c.userRepository = repository.NewPostgresUserRepository(db)
	c.activityRepository = activityRepo.NewActivityRepository(db)

	c.tokenService = token.NewTokenService()

	// Initialize use cases
	c.userUseCase = usecase.NewUserUseCase(c.userRepository, c.jwtManager, c.tokenService)
	c.activityUseCase = activityUseCase.NewActivityUseCase(c.activityRepository)

	// Initialize handlers
	c.userHandler = api.NewUserHandler(c.userUseCase)
	c.activityHandler = api.NewActivityHandler(c.activityUseCase)

	// Initialize middleware
	c.securityMiddleware = middleware.NewSecurityMiddleware(middleware.SecurityConfig{
		JWTManager: c.jwtManager,
	})
	c.errorMiddleware = middleware.NewErrorMiddleware()

	return c
}

// GetUserHandler returns the user handler
func (c *Container) GetUserHandler() *api.UserHandler {
	return c.userHandler
}

// GetActivityHandler returns the activity handler
func (c *Container) GetActivityHandler() *api.ActivityHandler {
	return c.activityHandler
}

// GetSecurityMiddleware returns the security middleware
func (c *Container) GetSecurityMiddleware() *middleware.SecurityMiddleware {
	return c.securityMiddleware
}

// GetErrorMiddleware returns the error middleware
func (c *Container) GetErrorMiddleware() *middleware.ErrorMiddleware {
	return c.errorMiddleware
}

// Close closes any resources held by the container
func (c *Container) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
