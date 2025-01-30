package usecase

import (
	"context"
	"dailyalu-server/internal/module/activity/domain"
)

type IActivityUseCase interface {
	Create(ctx context.Context, req *domain.CreateActivityRequest) (*domain.Activity, error)
	GetByID(ctx context.Context, id int) (*domain.Activity, error)
	Update(ctx context.Context, req *domain.UpdateActivityRequest) (*domain.Activity, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, req *domain.SearchActivityRequest) (*domain.ActivityResponse, error)
}