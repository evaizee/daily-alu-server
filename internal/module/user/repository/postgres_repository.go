package repository

import (
	"dailyalu-server/internal/module/user/domain"
	"database/sql"
	"fmt"
	"time"
)

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
		INSERT INTO users (id, email, name, password_hash, status, email_verification_token, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.PasswordHash, 
		user.Status, user.EmailVerificationToken, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) GetByID(id string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, status, email_verification_token, role, last_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, 
		&user.Status, &user.EmailVerificationToken, &user.Role,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *postgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	fmt.Println("pears")
	query := `
		SELECT id, email, name, password_hash, status, email_verification_token, role, last_login, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, 
		&user.Status, &user.EmailVerificationToken, &user.Role,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *postgresUserRepository) GetByVerificationToken(token string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, name, password_hash, status, email_verification_token, role, last_login, created_at, updated_at
		FROM users
		WHERE email_verification_token = $1
	`
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, 
		&user.Status, &user.EmailVerificationToken, &user.Role,
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
		SET email = $2, name = $3, status = $4, email_verification_token = $5, role = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, 
		user.Status, user.EmailVerificationToken, user.Role, user.UpdatedAt)
	return err
}

func (r *postgresUserRepository) UpdateLastLogin(id string, lastLogin time.Time) error {
	query := `
		UPDATE users
		SET last_login = $2, updated_at = $2
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, lastLogin)
	return err
}

func (r *postgresUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
