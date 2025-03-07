package repository

import (
	"context"
	"dailyalu-server/internal/module/activity/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type activityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) IActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	query := `
		INSERT INTO activities (user_id, child_id, type, details, happens_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		activity.UserID,
		activity.ChildID,
		activity.Type,
		activity.Details,
		activity.HappensAt,
		activity.CreatedAt,
		activity.UpdatedAt,
	).Scan(&activity.ID)

	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

func (r *activityRepository) GetByID(ctx context.Context, id int) (*domain.Activity, error) {
	query := `
		SELECT id, user_id, child_id, type, details, happens_at, created_at, updated_at
		FROM activities
		WHERE id = $1
	`
	activity := &domain.Activity{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.ChildID,
		&activity.Type,
		&activity.Details,
		&activity.HappensAt,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return activity, nil
}

func (r *activityRepository) Update(ctx context.Context, activity *domain.Activity) error {
	fmt.Println(activity)
	query := `
		UPDATE activities
		SET details = $1, happens_at = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		activity.Details,
		activity.HappensAt,
		activity.UpdatedAt,
		activity.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("activity not found or unauthorized")
	}

	return nil
}

func (r *activityRepository) Delete(ctx context.Context, id int) error {
	query := "UPDATE activities SET status = 20 WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("activity not found")
	}

	return nil
}

func (r *activityRepository) Search(ctx context.Context, req *domain.SearchActivityRequest) (*domain.ActivityResponse, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	argCount := 1

	if req.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, req.UserID)
		argCount++
	}

	if req.ChildID != 0 {
		conditions = append(conditions, fmt.Sprintf("child_id = $%d", argCount))
		args = append(args, req.ChildID)
		argCount++
	}

	if req.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, req.Type)
		argCount++
	}

	if !req.StartDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("happens_at >= $%d", argCount))
		args = append(args, req.StartDate)
		argCount++
	}

	if !req.EndDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("happens_at <= $%d", argCount))
		args = append(args, req.EndDate)
		argCount++
	}

	// Handle JSONB search
	if len(req.Details) > 0 {
		for key, value := range req.Details {
			jsonbCond := fmt.Sprintf("details->>'%s' = $%d", key, argCount)
			conditions = append(conditions, jsonbCond)
			valueStr, _ := json.Marshal(value)
			args = append(args, string(valueStr))
			argCount++
		}
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM activities WHERE %s
	`, strings.Join(conditions, " AND "))

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count activities: %w", err)
	}

	// Calculate pagination
	offset := (req.Page - 1) * req.PageSize
	totalPages := (int(total) + req.PageSize - 1) / req.PageSize

	// Get paginated records
	query := fmt.Sprintf(`
		SELECT id, user_id, child_id, type, details, happens_at, created_at, updated_at
		FROM activities
		WHERE %s
		ORDER BY happens_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), argCount, argCount+1)

	args = append(args, req.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search activities: %w", err)
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var activity domain.Activity
		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ChildID,
			&activity.Type,
			&activity.Details,
			&activity.HappensAt,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return &domain.ActivityResponse{
		Activities: activities,
		Pagination: domain.Pagination{
			Total:       total,
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			TotalPages:  totalPages,
		},
	}, nil
}
