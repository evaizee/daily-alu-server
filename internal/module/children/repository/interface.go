package repository

import (
	"dailyalu-server/internal/module/children/domain"
)

// IChildrenRepository defines the interface for children data access
type IChildrenRepository interface {
	Create(child *domain.Child) error
	GetByID(id int64) (*domain.Child, error)
	GetByUserID(userID string, page, pageSize int) ([]domain.Child, int64, error)
	Update(child *domain.Child) error
}
