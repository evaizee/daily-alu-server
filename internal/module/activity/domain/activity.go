package domain

import (
	"encoding/json"
	"time"
)

// Activity represents a baby activity record
type Activity struct {
	ID        int             `json:"id"`
	UserID    string          `json:"user_id"`
	ChildID   int             `json:"child_id"`
	Type      string          `json:"type"`
	Details   json.RawMessage `json:"details"`
	HappensAt time.Time       `json:"happens_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// CreateActivityRequest represents the request to create a new activity
type CreateActivityRequest struct {
	UserID    string          `json:"user_id" validate:"required"`
	ChildID   int             `json:"child_id"`
	Type      string          `json:"type" validate:"required"`
	Details   json.RawMessage `json:"details" validate:"required"`
	HappensAt string          `json:"happens_at" validate:"required"`
}

// UpdateActivityRequest represents the request to update an activity
type UpdateActivityRequest struct {
	ID        int             `json:"id" validate:"required"`
	UserID    string          `json:"user_id" validate:"required"`
	ChildID   string          `json:"child_id" validate:"required"`
	Details   json.RawMessage `json:"details" validate:"required"`
	HappensAt string          `json:"happens_at" validate:"required"`
}

// SearchActivityRequest represents the request to search activities
type SearchActivityRequest struct {
	UserID    string                 `json:"user_id"`
	ChildID   int                    `json:"child_id"`
	Type      string                 `json:"type"`
	StartDate time.Time              `json:"start_date"`
	EndDate   time.Time              `json:"end_date"`
	Details   map[string]interface{} `json:"details"` // For JSONB search
	Page      int                    `json:"page" validate:"min=1"`
	PageSize  int                    `json:"page_size" validate:"min=1,max=100"`
}

// Pagination represents pagination information
type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
}

// ActivityResponse represents the paginated response of activities
type ActivityResponse struct {
	Activities []Activity `json:"activities"`
	Pagination Pagination `json:"pagination"`
}
