package repository

import (
	"database/sql"
	"dailyalu-server/internal/module/user/domain"
	"time"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	UpdateLastLogin(id string, lastLogin time.Time) error
}

type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) IUserRepository {
	return &postgresUserRepository{db: db}
}

// Implementation of UserRepository interface
func (r *postgresUserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.PasswordHash, user.Role,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) GetByID(id string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, role, last_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *postgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, role, last_login, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *postgresUserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, role = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.Role, user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postgresUserRepository) UpdateLastLogin(id string, lastLogin time.Time) error {
	query := `UPDATE users SET last_login = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, lastLogin)
	return err
}
