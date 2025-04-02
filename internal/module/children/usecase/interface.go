package usecase

import (
	"dailyalu-server/internal/module/children/domain"
)

// IChildrenUseCase defines the interface for children business logic
type IChildrenUseCase interface {
	CreateChild(req *domain.CreateChildRequest) (*domain.Child, error)
	GetChild(id int64, userID string) (*domain.Child, error)
	GetChildren(req *domain.GetChildrenRequest) (*domain.ChildrenResponse, error)
	UpdateChild(req *domain.UpdateChildRequest) (*domain.Child, error)
}
