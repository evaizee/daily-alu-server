package api

import (
	"dailyalu-server/internal/module/user/domain"
	"dailyalu-server/internal/module/user/usecase"
	"dailyalu-server/internal/validator"
	"dailyalu-server/pkg/app_errors"
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
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
  }

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	_, err := h.userUseCase.Register(c.Context(), req)

	if err != nil {
		fmt.Println(err)
		errorStatus := fiber.StatusInternalServerError
		if err == usecase.ErrUserAlreadyExists {
			errorStatus = fiber.StatusConflict
		}

		return c.Status(errorStatus).JSON(fiber.Map{
			"code": errorStatus,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Your account is almost ready! To unlock all of our features, please verify your email address.",
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := &domain.LoginRequest{}

	if err := c.BodyParser(req); err != nil {
		fmt.Println("error = ",err)
		return app_errors.NewBadRequestError(err.Error())
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	loginResult, err := h.userUseCase.Login(req)

	if err == usecase.ErrInvalidCredentials {
		return app_errors.NewValidationError(err.Error())
	} else if err != nil {
		return app_errors.NewInternalError(err)
	}

	return c.JSON(loginResult)
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
	
	request :=  &domain.UpdateUserRequest{}
	if err := validator.ValidateRequest(c,request); err != nil {
		return err
	}

	request.ID = c.Params("id")

	user, err := h.userUseCase.UpdateUser(request)

	if err == usecase.ErrUserNotFound {
		return app_errors.NewError(app_errors.ErrCodeNotFound, err.Error())
	} else if err != nil {
		return app_errors.NewInternalError(err)
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
		return app_errors.NewValidationError("Verification token is required").AddMetadata("token",token)
	}

	if err := h.userUseCase.VerifyEmail(c.Context(), token); err != nil {
		return app_errors.NewBadRequestError(err.Error()).WithInternal(err).AddMetadata("token",token)
	}

	return c.JSON(fiber.Map{
			"message": "Email verified successfully",
	})
}

func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	request := &domain.UpdatePasswordRequest{}

	if err := validator.ValidateRequest(c,request); err != nil {
		return err
	}

	request.ID = c.Params("id")

	fmt.Println(request)

	err := h.userUseCase.UpdatePassword(request)

	if err == usecase.ErrUserNotFound {
		return app_errors.NewValidationError(err.Error())
	} else if err != nil {
		return app_errors.NewValidationError("Internal Server Error")
	}

	return c.JSON(fiber.Map{
		"message": "Password Updated successfully",
	})
}



// RefreshToken handles token refresh requests
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
    req := &domain.RefreshTokenRequest{}
		
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
