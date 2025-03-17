package response

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// PaginationRequest represents pagination parameters from client
type PaginationRequest struct {
	Page     int
	PageSize int
	SortBy   string
	SortDesc bool
}

// NewPagination creates a new Pagination instance
func NewPagination(totalItems int64, pageSize, currentPage int) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	
	hasNext := currentPage < totalPages
	hasPrevious := currentPage > 1
	
	var nextPage *int
	var previousPage *int
	
	if hasNext {
		next := currentPage + 1
		nextPage = &next
	}
	
	if hasPrevious {
		prev := currentPage - 1
		previousPage = &prev
	}
	
	return Pagination{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}
}

// SuccessWithPagination returns a standardized paginated success response
func SuccessWithPagination(c *fiber.Ctx, statusCode int, message string, data interface{}, pagination Pagination) error {
	return c.Status(statusCode).JSON(PaginatedResponse{
		Code:       statusCode,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ParsePaginationRequest parses pagination parameters from the request
func ParsePaginationRequest(c *fiber.Ctx) PaginationRequest {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	
	pageSize := c.QueryInt("page_size", 10)
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}
	
	sortBy := c.Query("sort", "")
	sortDesc := c.Query("order", "asc") == "desc"
	
	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
}
