package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/postgres"
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

func (r *UserRepo) Save(ctx context.Context, u *domain.User, cred *domain.UserCredential, profile *domain.Profile) error {
	executor := postgres.GetExecutor(ctx, r.db)

	// 1. Insert into users
	queryUser := `INSERT INTO users (username, username_updated_at, status, role, created_at, updated_at)
	              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := executor.QueryRowContext(ctx, queryUser,
		u.Username, u.UsernameUpdatedAt, u.Status, u.Role, u.CreatedAt, u.UpdatedAt,
	).Scan(&u.ID)
	if err != nil {
		return err
	}

	// 2. Set user ID in dependent structs
	cred.UserID = u.ID
	profile.UserID = u.ID

	// 3. Insert into user_credentials
	queryCred := `INSERT INTO user_credentials (user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at)
	              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = executor.QueryRowContext(ctx, queryCred,
		cred.UserID, cred.Type, cred.Identifier, cred.SecretHash, cred.IsVerified, cred.IsPrimary, cred.CreatedAt, cred.UpdatedAt,
	).Scan(&cred.ID)
	if err != nil {
		return err
	}

	// 4. Insert into profiles
	queryProfile := `INSERT INTO profiles (user_id, full_name, avatar_url, gender, dob, updated_at)
	                 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = executor.ExecContext(ctx, queryProfile,
		profile.UserID, profile.FullName, profile.AvatarURL, string(profile.Gender), profile.DOB, profile.UpdatedAt,
	)

	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT id, username, username_updated_at, status, role, created_at, updated_at
	          FROM users WHERE id = $1`
	err := executor.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.UsernameUpdatedAt, &u.Status, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT id, username, username_updated_at, status, role, created_at, updated_at
	          FROM users WHERE username = $1`
	err := executor.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.UsernameUpdatedAt, &u.Status, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByIdentifier(ctx context.Context, identifier string) (*domain.User, error) {
	var u domain.User
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT u.id, u.username, u.username_updated_at, u.status, u.role, u.created_at, u.updated_at
			FROM users u
			LEFT JOIN user_credentials uc 
			ON u.id = uc.user_id 
			WHERE uc.identifier = $1`
	err := executor.QueryRowContext(ctx, query, identifier).Scan(
		&u.ID, &u.Username, &u.UsernameUpdatedAt, &u.Status, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := executor.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) ExistsByIdentifier(ctx context.Context, identifier string) (bool, error) {
	var exists bool
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT EXISTS(SELECT 1 FROM user_credentials WHERE identifier = $1)`
	err := executor.QueryRowContext(ctx, query, identifier).Scan(&exists)
	return exists, err
}

func (r *UserRepo) FindCredentialByIdentifier(ctx context.Context, identifier string) (*domain.UserCredential, error) {
	var c domain.UserCredential
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT id, user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at
	          FROM user_credentials WHERE identifier = $1`
	err := executor.QueryRowContext(ctx, query, identifier).Scan(
		&c.ID, &c.UserID, &c.Type, &c.Identifier, &c.SecretHash, &c.IsVerified, &c.IsPrimary, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *UserRepo) FindCredentialsByUserID(ctx context.Context, userID string) ([]*domain.UserCredential, error) {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT id, user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at
	          FROM user_credentials WHERE user_id = $1`
	rows, err := executor.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []*domain.UserCredential
	for rows.Next() {
		var c domain.UserCredential
		err := rows.Scan(
			&c.ID, &c.UserID, &c.Type, &c.Identifier, &c.SecretHash, &c.IsVerified, &c.IsPrimary, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		creds = append(creds, &c)
	}
	return creds, nil
}

func (r *UserRepo) FindProfileByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	var p domain.Profile
	var genderStr string
	executor := postgres.GetExecutor(ctx, r.db)
	query := `SELECT user_id, full_name, avatar_url, gender, dob, updated_at
	          FROM profiles WHERE user_id = $1`
	err := executor.QueryRowContext(ctx, query, userID).Scan(
		&p.UserID, &p.FullName, &p.AvatarURL, &genderStr, &p.DOB, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	gender := domain.Gender(genderStr)
	if !gender.IsValid() {
		return nil, fmt.Errorf("invalid gender value from database: %s", genderStr)
	}
	p.Gender = gender

	return &p, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, p *domain.Profile) error {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `UPDATE profiles SET full_name = $1, avatar_url = $2, gender = $3, dob = $4, updated_at = $5
	          WHERE user_id = $6`
	_, err := executor.ExecContext(ctx, query, p.FullName, p.AvatarURL, string(p.Gender), p.DOB, p.UpdatedAt, p.UserID)
	return err
}

func (r *UserRepo) UpdateCredential(ctx context.Context, c *domain.UserCredential) error {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `UPDATE user_credentials SET identifier = $1, secret_hash = $2, is_verified = $3, is_primary = $4, updated_at = $5
	          WHERE id = $6`
	_, err := executor.ExecContext(ctx, query, c.Identifier, c.SecretHash, c.IsVerified, c.IsPrimary, c.UpdatedAt, c.ID)
	return err
}

func (r *UserRepo) SaveCredential(ctx context.Context, c *domain.UserCredential) error {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `INSERT INTO user_credentials (user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := executor.QueryRowContext(ctx, query,
		c.UserID, c.Type, c.Identifier, c.SecretHash, c.IsVerified, c.IsPrimary, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.ID)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `UPDATE user_credentials SET secret_hash = $1, updated_at = NOW()
	          WHERE user_id = $2 AND is_primary = true`
	_, err := executor.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	executor := postgres.GetExecutor(ctx, r.db)
	query := `UPDATE users SET username = $1, username_updated_at = $2, status = $3, role = $4, updated_at = $5
	          WHERE id = $6`
	_, err := executor.ExecContext(ctx, query, u.Username, u.UsernameUpdatedAt, u.Status, u.Role, u.UpdatedAt, u.ID)
	return err
}
