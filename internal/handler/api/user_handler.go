package api

import (
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/validator"
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
        return c.SendStatus(500)
    }

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	user, err := h.userUseCase.Register(c.Context(), req)
	
	if err == usecase.ErrUserAlreadyExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := &domain.LoginRequest{}

	if err := c.BodyParser(req); err != nil {
        fmt.Println("error = ",err)
        return c.SendStatus(500)
    }

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	accessToken, refreshToken, err := h.userUseCase.Login(req)
	if err == usecase.ErrInvalidCredentials {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userUseCase.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var input struct {
		Email string `json:"email" validate:"required,email"`
		Name  string `json:"name" validate:"required"`
	}

	if err := validator.ValidateRequest(c, &input); err != nil {
		return err
	}

	user, err := h.userUseCase.UpdateUser(id, input.Email, input.Name)
	if err == usecase.ErrUserNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userUseCase.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
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
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Verification token is required",
        })
    }

    if err := h.userUseCase.VerifyEmail(c.Context(), token); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Email verified successfully",
    })
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshToken handles token refresh requests
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
    req := &RefreshTokenRequest{}
    if err := c.BodyParser(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := validator.ValidateStruct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   "Validation failed",
            "details": err,
        })
    }

    // Generate new token pair using refresh token
    accessToken, newRefreshToken, err := h.userUseCase.RefreshToken(req.RefreshToken)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid or expired refresh token",
        })
    }

    return c.JSON(fiber.Map{
        "access_token": accessToken,
        "refresh_token": newRefreshToken,
    })
}
