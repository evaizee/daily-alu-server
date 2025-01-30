package repository

import (
	"context"
	"dailyalu-server/internal/module/activity/domain"
)

type IActivityRepository interface {
	Create(ctx context.Context, activity *domain.Activity) error
	GetByID(ctx context.Context, id int) (*domain.Activity, error)
	Update(ctx context.Context, activity *domain.Activity) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, req *domain.SearchActivityRequest) (*domain.ActivityResponse, error)
}