package usecase

import (
	"context"
	"dailyalu-server/internal/module/activity/domain"
	"dailyalu-server/internal/module/activity/repository"
	"dailyalu-server/internal/utils"
	"fmt"
	"time"
)

type activityUseCase struct {
	repo repository.IActivityRepository
}

func NewActivityUseCase(repo repository.IActivityRepository) IActivityUseCase {
	return &activityUseCase{
		repo: repo,
	}
}

func (uc *activityUseCase) Create(ctx context.Context, req *domain.CreateActivityRequest) (*domain.Activity, error) {
	now := time.Now()

	happensAt, err := utils.TimeLocationParsing(ctx, req.HappensAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	activity := &domain.Activity{
		UserID:    req.UserID,
		ChildID:   req.ChildID,
		Type:      req.Type,
		Details:   req.Details,
		HappensAt: happensAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err = uc.repo.Create(ctx, activity); err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return activity, nil
}

func (uc *activityUseCase) GetByID(ctx context.Context, id int) (*domain.Activity, error) {
	activity, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("activity not found")
	}
	return activity, nil
}

func (uc *activityUseCase) Update(ctx context.Context, req *domain.UpdateActivityRequest) (*domain.Activity, error) {
	activity, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("activity not found")
	}

	// Verify ownership
	if activity.UserID != req.UserID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update fields
	activity.Details = req.Details
	activity.UpdatedAt = time.Now()
	activity.HappensAt, err = utils.TimeLocationParsing(ctx, req.HappensAt)
	if err != nil {	
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	if err := uc.repo.Update(ctx, activity); err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	return activity, nil
}

func (uc *activityUseCase) Delete(ctx context.Context, id int) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}
	return nil
}

func (uc *activityUseCase) Search(ctx context.Context, req *domain.SearchActivityRequest) (*domain.ActivityResponse, error) {
	// Set default pagination values if not provided
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	response, err := uc.repo.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search activities: %w", err)
	}

	return response, nil
}
