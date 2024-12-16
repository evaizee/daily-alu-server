package handler

import (
	"dailyalu-server/internal/module/user/usecase"

	"github.com/gofiber/fiber/v2"
)

type verifyHandler struct {
	userUseCase usecase.IUserUseCase
}

func NewVerifyHandler(userUseCase usecase.IUserUseCase) *verifyHandler {
	return &verifyHandler{
		userUseCase: userUseCase,
	}
}

// VerifyEmail handles email verification
// @Summary Verify email address
// @Description Verifies a user's email address using the verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/verify-email [get]
func (h *verifyHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verification token is required",
		})
	}

	if err := h.userUseCase.VerifyEmail(c.Context(), token); err != nil {
		switch err {
		case usecase.ErrInvalidVerificationToken:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid verification token",
			})
		case usecase.ErrVerificationTokenExpired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Verification token has expired",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify email",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
	})
}
