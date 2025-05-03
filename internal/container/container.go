package container

import (
	"dailyalu-server/internal/handler/api"
	"dailyalu-server/internal/middleware"
	activityRepo "dailyalu-server/internal/module/activity/repository"
	activityUseCase "dailyalu-server/internal/module/activity/usecase"
	childrenRepo "dailyalu-server/internal/module/children/repository"
	childrenUseCase "dailyalu-server/internal/module/children/usecase"
	"dailyalu-server/internal/module/user/repository"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/security/token"
	"dailyalu-server/internal/service/mailer"
	"database/sql"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
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
	childrenRepository childrenRepo.IChildrenRepository

	// Use Cases
	userUseCase     usecase.IUserUseCase
	activityUseCase activityUseCase.IActivityUseCase
	childrenUseCase childrenUseCase.IChildrenUseCase

	// Handlers
	userHandler     *api.UserHandler
	activityHandler *api.ActivityHandler
	childrenHandler *api.ChildrenHandler

	// Middleware
	securityMiddleware *middleware.SecurityMiddleware
	errorMiddleware    *middleware.ErrorMiddleware

	//Token Verification Service
	tokenService *token.TokenService

	//Mailer Service
	mailerService *mailer.MailerService
}

// NewContainer creates a new dependency injection container
func NewContainer(db *sql.DB, newSes *sesv2.Client, jwtSecret, jwtRefreshSecretKey string, jwtExpiry, jwtRefreshExpiry time.Duration) *Container {
	c := &Container{
		db: db,
	}

	// Initialize JWT manager
	c.jwtManager = jwt.NewJWTManager(jwtSecret, jwtRefreshSecretKey, jwtExpiry, jwtRefreshExpiry)

	// Initialize mailer
	c.mailerService = mailer.NewMailerService(newSes)

	// Initialize repositories
	c.userRepository = repository.NewPostgresUserRepository(db)
	c.activityRepository = activityRepo.NewActivityRepository(db)
	c.childrenRepository = childrenRepo.NewPostgresChildrenRepository(db)

	c.tokenService = token.NewTokenService()

	// Initialize use cases
	c.userUseCase = usecase.NewUserUseCase(c.userRepository, c.jwtManager, c.tokenService, c.mailerService)
	c.activityUseCase = activityUseCase.NewActivityUseCase(c.activityRepository)
	c.childrenUseCase = childrenUseCase.NewChildrenUseCase(c.childrenRepository)

	// Initialize handlers
	c.userHandler = api.NewUserHandler(c.userUseCase)
	c.activityHandler = api.NewActivityHandler(c.activityUseCase)
	c.childrenHandler = api.NewChildrenHandler(c.childrenUseCase)

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

// GetChildrenHandler returns the children handler
func (c *Container) GetChildrenHandler() *api.ChildrenHandler {
	return c.childrenHandler
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
