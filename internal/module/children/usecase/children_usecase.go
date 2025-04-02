package usecase

import (
	"dailyalu-server/internal/module/children/domain"
	"dailyalu-server/internal/module/children/repository"
	"database/sql"
	"math"
)

// ChildrenUseCase implements the children business logic
type ChildrenUseCase struct {
	childrenRepo repository.IChildrenRepository
}

// NewChildrenUseCase creates a new children use case
func NewChildrenUseCase(childrenRepo repository.IChildrenRepository) IChildrenUseCase {
	return &ChildrenUseCase{
		childrenRepo: childrenRepo,
	}
}

// CreateChild creates a new child
func (u *ChildrenUseCase) CreateChild(req *domain.CreateChildRequest) (*domain.Child, error) {
	// Create child entity
	child := &domain.Child{
		UserID:  req.UserID,
		Name:    req.Name,
		Details: req.Details,
	}

	// Save to repository
	err := u.childrenRepo.Create(child)
	if err != nil {
		return nil, err
	}

	return child, nil
}

// GetChild retrieves a child by ID and validates ownership
func (u *ChildrenUseCase) GetChild(id int64, userID string) (*domain.Child, error) {
	// Get child from repository
	child, err := u.childrenRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if child exists
	if child == nil {
		return nil, ErrChildNotFound
	}

	// Validate ownership
	if child.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	return child, nil
}

// GetChildren retrieves children for a user with pagination
func (u *ChildrenUseCase) GetChildren(req *domain.GetChildrenRequest) (*domain.ChildrenResponse, error) {
	// Get children from repository with pagination
	children, total, err := u.childrenRepo.GetByUserID(req.UserID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	// Create response
	response := &domain.ChildrenResponse{
		Children: children,
		Pagination: domain.Pagination{
			Total:       total,
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			TotalPages:  totalPages,
		},
	}

	return response, nil
}

// UpdateChild updates an existing child
func (u *ChildrenUseCase) UpdateChild(req *domain.UpdateChildRequest) (*domain.Child, error) {
	// Check if child exists and belongs to the user
	child, err := u.childrenRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if child == nil {
		return nil, ErrChildNotFound
	}

	// Validate ownership
	if child.UserID != req.UserID {
		return nil, ErrUnauthorizedAccess
	}

	// Update child entity
	child.Name = req.Name
	child.Details = req.Details

	// Save to repository
	err = u.childrenRepo.Update(child)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrChildNotFound
		}
		return nil, err
	}

	return child, nil
}
