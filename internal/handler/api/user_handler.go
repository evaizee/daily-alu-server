package api

import (
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/utils"
	"dailyalu-server/internal/validator"
	"dailyalu-server/pkg/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUseCase usecase.IUserUseCase
}

func NewUserHandler(userUseCase usecase.IUserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	req := &domain.RegisterRequest{}
	if err := c.BodyParser(req); err != nil {
		fmt.Println("error = ",err)
		return response.NewBadRequestError("Invalid request body")
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	user, err := h.userUseCase.Register(c.Context(), req)

	if err != nil {
		fmt.Println(err)
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusCreated, 
		"Your account is almost ready! To unlock all of our features, please verify your email address.", 
		user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := &domain.LoginRequest{}

	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	loginResult, err := h.userUseCase.Login(req)

	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "Login successful", loginResult)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID := utils.GetUserIDFromContext(c)
	
	user, err := h.userUseCase.GetUser(userID)

	if err != nil {
		return response.MapDomainError(err)
	}

	if user == nil {
		return response.NewNotFoundError("User not found")
	}

	return response.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	request := &domain.UpdateUserRequest{}
	if err := validator.ValidateRequest(c, request); err != nil {
		return err
	}

	// Get user ID from JWT token
	request.ID = utils.GetUserIDFromContext(c)

	user, err := h.userUseCase.UpdateUser(request)

	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userUseCase.DeleteUser(id); err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
    // Check query parameter first
	token := c.Query("token")
	
	// If not in query, check path parameter
	if token == "" {
		token = c.Params("token")
	}

	if token == "" {
		return response.NewValidationError("Verification token is required").AddMetadata("token", token)
	}

	if err := h.userUseCase.VerifyEmail(c.Context(), token); err != nil {
		return response.NewBadRequestError(err.Error()).WithInternal(err).AddMetadata("token", token)
	}

	return response.Success(c, fiber.StatusOK, "Email verified successfully", nil)
}

func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	request := &domain.UpdatePasswordRequest{}

	if err := validator.ValidateRequest(c, request); err != nil {
		return err
	}

	// Get user ID from JWT token
	request.ID = utils.GetUserIDFromContext(c)

	err := h.userUseCase.UpdatePassword(request)

	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "Password updated successfully", nil)
}

// RefreshToken handles token refresh requests
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	req := &domain.RefreshTokenRequest{}
	
	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		return response.NewValidationErrorWithDetails("Validation failed", errors)
	}

	// Generate new token pair using refresh token
	accessToken, newRefreshToken, err := h.userUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		return response.NewUnauthorizedError("Invalid or expired refresh token")
	}

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token": accessToken,
		"refresh_token": newRefreshToken,
	})
}

// ForgotPassword handles password reset requests
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	req := &domain.ForgotPasswordRequest{}
	
	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	if err := h.userUseCase.ForgotPassword(req); err != nil {
		// Don't expose detailed errors to avoid email enumeration
		// Just log the error internally
		fmt.Println("Error in forgot password:", err)
		// But still map domain errors for proper handling
		return response.MapDomainError(err)
	}

	// Always return success even if email doesn't exist (for security)
	return response.Success(
		c, 
		fiber.StatusOK, 
		"If your email is registered with us, you will receive password reset instructions shortly", 
		nil,
	)
}

// ResetPassword handles password reset with token
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	req := &domain.ResetPasswordRequest{}
	
	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// If token is not in the request body, check query parameter
	if req.Token == "" {
		req.Token = c.Query("token")
	}

	if req.Token == "" {
		return response.NewValidationError("Reset token is required")
	}

	if err := h.userUseCase.ResetPassword(req); err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(
		c, 
		fiber.StatusOK, 
		"Your password has been reset successfully", 
		nil,
	)
}

// AdminGetUser allows admins to get any user by ID
func (h *UserHandler) AdminGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userUseCase.GetUser(id)

	if err != nil {
		return response.MapDomainError(err)
	}

	if user == nil {
		return response.NewNotFoundError("User not found")
	}

	return response.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

// AdminUpdateUser allows admins to update any user by ID
func (h *UserHandler) AdminUpdateUser(c *fiber.Ctx) error {
	request := &domain.UpdateUserRequest{}
	if err := validator.ValidateRequest(c, request); err != nil {
		return err
	}

	request.ID = c.Params("id")

	user, err := h.userUseCase.UpdateUser(request)

	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "User updated successfully", user)
}
