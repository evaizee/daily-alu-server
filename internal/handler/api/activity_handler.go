package api

import (
	"dailyalu-server/internal/module/activity/domain"
	"dailyalu-server/internal/module/activity/usecase"
	"dailyalu-server/internal/utils"
	"dailyalu-server/internal/validator"
	"dailyalu-server/pkg/response"
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
	userID := utils.GetUserIDFromContext(c)
	req := &domain.CreateActivityRequest{}

	if err := c.BodyParser(req); err != nil {
		fmt.Println("error = ", err)
		return response.NewBadRequestError("Invalid request body")
	}

	req.UserID = userID

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	activity, err := h.activityUseCase.Create(c.Context(), req)
	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusCreated, "Activity created successfully", activity)
}

func (h *ActivityHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.NewBadRequestError("Invalid activity ID format")
	}

	activity, err := h.activityUseCase.GetByID(c.Context(), id)
	if err != nil {
		return response.NewNotFoundError("Activity not found")
	}

	return response.Success(c, fiber.StatusOK, "Activity retrieved successfully", activity)
}

func (h *ActivityHandler) Update(c *fiber.Ctx) error {
	userID := utils.GetUserIDFromContext(c)

	req := &domain.UpdateActivityRequest{}
	req.UserID = userID
	
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err.Error())
		return response.NewBadRequestError("Invalid request body")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fmt.Println(err.Error())
		return response.NewBadRequestError("Invalid activity ID format")
	}

	req.ID = id

	if err := validator.ValidateRequest(c, req); err != nil {
		fmt.Println(err.Error())
		return err
	}

	activity, err := h.activityUseCase.Update(c.Context(), req)
	if err != nil {
		fmt.Println(err.Error())
		return response.MapDomainError(err)
	}
	
	return response.Success(c, fiber.StatusOK, "Activity updated successfully", activity)
}

func (h *ActivityHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.NewBadRequestError("Invalid activity ID format")
	}

	if err := h.activityUseCase.Delete(c.Context(), id); err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "Activity deleted successfully", nil)
}

func (h *ActivityHandler) Search(c *fiber.Ctx) error {
	userID := utils.GetUserIDFromContext(c)
	req := &domain.SearchActivityRequest{
		UserID:   userID,
		Type:     c.Query("type"),
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 10),
	}

	// Parse dates if provided
	if startDate := c.Query("start_date"); startDate != "" {
		date, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return response.NewBadRequestError("Invalid start_date format")
		}
		req.StartDate = date
	}

	if endDate := c.Query("end_date"); endDate != "" {
		date, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return response.NewBadRequestError("Invalid end_date format")
		}
		req.EndDate = date
	}

	// Parse details if provided
	if details := c.Query("details"); details != "" {
		var detailsMap map[string]interface{}
		if err := json.Unmarshal([]byte(details), &detailsMap); err != nil {
			return response.NewBadRequestError("Invalid details format")
		}
		req.Details = detailsMap
	}

	// Get search results
	activityResponse, err := h.activityUseCase.Search(c.Context(), req)
	if err != nil {
		return response.MapDomainError(err)
	}

	// Create pagination from the response
	pagination := response.NewPagination(
		activityResponse.Pagination.Total, 
		activityResponse.Pagination.PageSize, 
		activityResponse.Pagination.CurrentPage,
	)

	// Return paginated response
	return response.SuccessWithPagination(
		c, 
		fiber.StatusOK, 
		"Activities retrieved successfully", 
		activityResponse.Activities, 
		pagination,
	)
}
