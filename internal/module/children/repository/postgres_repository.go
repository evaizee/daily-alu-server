package repository

import (
	"dailyalu-server/internal/module/children/domain"
	"database/sql"
	"time"
)

// PostgresChildrenRepository implements the children repository interface using PostgreSQL
type PostgresChildrenRepository struct {
	db *sql.DB
}

// NewPostgresChildrenRepository creates a new PostgreSQL children repository
func NewPostgresChildrenRepository(db *sql.DB) IChildrenRepository {
	return &PostgresChildrenRepository{
		db: db,
	}
}

// Create inserts a new child record into the database
func (r *PostgresChildrenRepository) Create(child *domain.Child) error {
	query := `
		INSERT INTO children (user_id, name, details, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	child.CreatedAt = now
	child.UpdatedAt = now

	var details []byte
	if child.Details != nil {
		details = child.Details
	}

	err := r.db.QueryRow(
		query,
		child.UserID,
		child.Name,
		details,
		child.CreatedAt,
		child.UpdatedAt,
	).Scan(&child.ID)

	return err
}

// GetByID retrieves a child by ID
func (r *PostgresChildrenRepository) GetByID(id int64) (*domain.Child, error) {
	query := `
		SELECT id, user_id, name, details, created_at, updated_at
		FROM children
		WHERE id = $1
	`

	var child domain.Child
	var details sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&child.ID,
		&child.UserID,
		&child.Name,
		&details,
		&child.CreatedAt,
		&child.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if details.Valid {
		child.Details = []byte(details.String)
	}

	return &child, nil
}

// GetByUserID retrieves children by user ID with pagination
func (r *PostgresChildrenRepository) GetByUserID(userID string, page, pageSize int) ([]domain.Child, int64, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM children WHERE user_id = $1`
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, user_id, name, details, created_at, updated_at
		FROM children
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var children []domain.Child
	for rows.Next() {
		var child domain.Child
		var details sql.NullString

		err := rows.Scan(
			&child.ID,
			&child.UserID,
			&child.Name,
			&details,
			&child.CreatedAt,
			&child.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if details.Valid {
			child.Details = []byte(details.String)
		}

		children = append(children, child)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return children, total, nil
}

// Update updates an existing child record
func (r *PostgresChildrenRepository) Update(child *domain.Child) error {
	query := `
		UPDATE children
		SET name = $1, details = $2, updated_at = $3
		WHERE id = $4 AND user_id = $5
	`

	child.UpdatedAt = time.Now()

	var details []byte
	if child.Details != nil {
		details = child.Details
	}

	result, err := r.db.Exec(
		query,
		child.Name,
		details,
		child.UpdatedAt,
		child.ID,
		child.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
