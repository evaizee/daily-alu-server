package api

import (
	"dailyalu-server/internal/module/activity/domain"
	"dailyalu-server/internal/module/activity/usecase"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/validator"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ActivityHandler struct {
	activityUseCase usecase.IActivityUseCase
}

func NewActivityHandler(activityUseCase usecase.IActivityUseCase) *ActivityHandler {
	return &ActivityHandler{
		activityUseCase: activityUseCase,
	}
}

func (h *ActivityHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID
	req := &domain.CreateActivityRequest{}

	if err := c.BodyParser(req); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(500)
	}

	req.UserID = userID

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	activity, err := h.activityUseCase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(activity)
}

func (h *ActivityHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to parse activity ID",
		})
	}

	activity, err := h.activityUseCase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Activity not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(activity)
}

func (h *ActivityHandler) Update(c *fiber.Ctx) error {
	fmt.Println("update activity")
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID

	req := &domain.UpdateActivityRequest{}
	req.UserID = userID
	
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to parse activity ID",
		})
	}

	req.ID = id

	if err := validator.ValidateRequest(c, req); err != nil {
		fmt.Println(err.Error())
		return err
	}

	activity, err := h.activityUseCase.Update(c.Context(), req)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
  fmt.Println("activity = ", activity)
	return c.Status(fiber.StatusCreated).JSON(activity)
}

func (h *ActivityHandler) Delete(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to parse activity ID",
		})
	}

	if err := h.activityUseCase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ActivityHandler) Search(c *fiber.Ctx) error {

	claims := c.Locals("user").(*jwt.Claims)
	req := &domain.SearchActivityRequest{
		UserID:   claims.UserID,
		Type:     c.Query("type"),
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 10),
	}

	// Parse dates if provided
	if startDate := c.Query("start_date"); startDate != "" {
		date, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start_date format",
			})
		}
		req.StartDate = date
	}

	if endDate := c.Query("end_date"); endDate != "" {
		date, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end_date format",
			})
		}
		req.EndDate = date
	}

	// Parse details if provided
	if details := c.Query("details"); details != "" {
		var detailsMap map[string]interface{}
		if err := json.Unmarshal([]byte(details), &detailsMap); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid details format",
			})
		}
		req.Details = detailsMap
	}

	response, err := h.activityUseCase.Search(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
