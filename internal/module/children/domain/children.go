package domain

import (
	"encoding/json"
	"time"
)

// Child represents a child record
type Child struct {
	ID        int64           `json:"id"`
	UserID    string          `json:"user_id"`
	Name      string          `json:"name"`
	Details   json.RawMessage `json:"details,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CreateChildRequest represents the request to create a new child
type CreateChildRequest struct {
	UserID  string          `json:"user_id" validate:"required"`
	Name    string          `json:"name" validate:"required"`
	Details json.RawMessage `json:"details,omitempty"`
}

// UpdateChildRequest represents the request to update a child
type UpdateChildRequest struct {
	ID      int64           `json:"id" validate:"required"`
	UserID  string          `json:"user_id" validate:"required"`
	Name    string          `json:"name" validate:"required"`
	Details json.RawMessage `json:"details,omitempty"`
}

// GetChildrenRequest represents the request to get children with pagination
type GetChildrenRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

// ChildrenResponse represents the paginated response of children
type ChildrenResponse struct {
	Children   []Child     `json:"children"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
}
