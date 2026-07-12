package postgres

import (
	"context"
	"database/sql"
	"errors"
	"user-service/internal/domain"
	"user-service/internal/repo"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

var _ repo.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Save(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (username, email, phone, password_hash, full_name, status, created_at, updated_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		u.Username, u.Email, u.Phone, u.PasswordHash, u.FullName, u.Status, u.CreatedAt, u.UpdatedAt,
	).Scan(&u.ID)
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, username, email, phone, password_hash, full_name, status, created_at, updated_at
	          FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.Phone, &u.PasswordHash, &u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, username, email, phone, password_hash, full_name, status, created_at, updated_at
	          FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.Phone, &u.PasswordHash, &u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}
