package api

import (
	"dailyalu-server/internal/module/children/domain"
	"dailyalu-server/internal/module/children/usecase"
	"dailyalu-server/internal/security/jwt"
	"dailyalu-server/internal/validator"
	"dailyalu-server/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ChildrenHandler handles HTTP requests for children
type ChildrenHandler struct {
	childrenUseCase usecase.IChildrenUseCase
}

// NewChildrenHandler creates a new children handler
func NewChildrenHandler(childrenUseCase usecase.IChildrenUseCase) *ChildrenHandler {
	return &ChildrenHandler{
		childrenUseCase: childrenUseCase,
	}
}

// CreateChild handles the creation of a new child
func (h *ChildrenHandler) CreateChild(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID
	if userID == "" {
		return response.NewUnauthorizedError("Authentication required")
	}

	// Parse request body
	req := &domain.CreateChildRequest{}
	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	// Set user ID from authenticated user
	req.UserID = userID

	// Validate request
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// Create child
	child, err := h.childrenUseCase.CreateChild(req)
	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusCreated, "Child created successfully", child)
}

// GetChild handles retrieving a child by ID
func (h *ChildrenHandler) GetChild(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID
	if userID == "" {
		return response.NewUnauthorizedError("Authentication required")
	}

	// Parse child ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return response.NewBadRequestError("Invalid child ID")
	}

	// Get child
	child, err := h.childrenUseCase.GetChild(id, userID)
	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "Child retrieved successfully", child)
}

// GetChildren handles retrieving all children for a user with pagination
func (h *ChildrenHandler) GetChildren(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID
	if userID == "" {
		return response.NewUnauthorizedError("Authentication required")
	}

	// Parse pagination parameters
	paginationReq := response.ParsePaginationRequest(c)

	// Create request
	req := &domain.GetChildrenRequest{
		UserID:   userID,
		Page:     paginationReq.Page,
		PageSize: paginationReq.PageSize,
	}

	// Get children
	result, err := h.childrenUseCase.GetChildren(req)
	if err != nil {
		return response.MapDomainError(err)
	}

	// Create pagination metadata
	pagination := response.NewPagination(
		result.Pagination.Total,
		result.Pagination.PageSize,
		result.Pagination.CurrentPage,
	)

	return response.SuccessWithPagination(
		c,
		fiber.StatusOK,
		"Children retrieved successfully",
		result.Children,
		pagination,
	)
}

// UpdateChild handles updating an existing child
func (h *ChildrenHandler) UpdateChild(c *fiber.Ctx) error {
	// Get authenticated user ID from context
	claims := c.Locals("user").(*jwt.Claims)
	userID := claims.UserID
	if userID == "" {
		return response.NewUnauthorizedError("Authentication required")
	}

	// Parse child ID from path parameter
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return response.NewBadRequestError("Invalid child ID")
	}

	// Parse request body
	req := &domain.UpdateChildRequest{}
	if err := c.BodyParser(req); err != nil {
		return response.NewBadRequestError("Invalid request body")
	}

	// Set ID and user ID
	req.ID = id
	req.UserID = userID

	// Validate request
	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	// Update child
	child, err := h.childrenUseCase.UpdateChild(req)
	if err != nil {
		return response.MapDomainError(err)
	}

	return response.Success(c, fiber.StatusOK, "Child updated successfully", child)
}