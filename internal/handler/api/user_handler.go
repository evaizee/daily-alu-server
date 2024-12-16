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
	var req *domain.LoginRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		return err
	}

	token, err := h.userUseCase.Login(req)
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
		"token": token,
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
