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

func (r *UserRepo) Save(ctx context.Context, u *domain.User, cred *domain.UserCredential, profile *domain.Profile) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert into users
	queryUser := `INSERT INTO users (username, username_updated_at, status, role, created_at, updated_at)
	              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tx.QueryRowContext(ctx, queryUser,
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
	err = tx.QueryRowContext(ctx, queryCred,
		cred.UserID, cred.Type, cred.Identifier, cred.SecretHash, cred.IsVerified, cred.IsPrimary, cred.CreatedAt, cred.UpdatedAt,
	).Scan(&cred.ID)
	if err != nil {
		return err
	}

	// 4. Insert into profiles
	queryProfile := `INSERT INTO profiles (user_id, full_name, avatar_url, gender, dob, updated_at)
	                 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, queryProfile,
		profile.UserID, profile.FullName, profile.AvatarURL, profile.Gender, profile.DOB, profile.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, username, username_updated_at, status, role, created_at, updated_at
	          FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
	query := `SELECT id, username, username_updated_at, status, role, created_at, updated_at
	          FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).Scan(
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
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) ExistsByIdentifier(ctx context.Context, identifier string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_credentials WHERE identifier = $1)`
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(&exists)
	return exists, err
}

func (r *UserRepo) FindCredentialByIdentifier(ctx context.Context, identifier string) (*domain.UserCredential, error) {
	var c domain.UserCredential
	query := `SELECT id, user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at
	          FROM user_credentials WHERE identifier = $1`
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
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
	query := `SELECT id, user_id, type, identifier, secret_hash, is_verified, is_primary, created_at, updated_at
	          FROM user_credentials WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
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
	query := `SELECT user_id, full_name, avatar_url, gender, dob, updated_at
	          FROM profiles WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&p.UserID, &p.FullName, &p.AvatarURL, &p.Gender, &p.DOB, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, p *domain.Profile) error {
	query := `UPDATE profiles SET full_name = $1, avatar_url = $2, gender = $3, dob = $4, updated_at = $5
	          WHERE user_id = $6`
	_, err := r.db.ExecContext(ctx, query, p.FullName, p.AvatarURL, p.Gender, p.DOB, p.UpdatedAt, p.UserID)
	return err
}

func (r *UserRepo) UpdateCredential(ctx context.Context, c *domain.UserCredential) error {
	query := `UPDATE user_credentials SET identifier = $1, secret_hash = $2, is_verified = $3, is_primary = $4, updated_at = $5
	          WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, c.Identifier, c.SecretHash, c.IsVerified, c.IsPrimary, c.UpdatedAt, c.ID)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	query := `UPDATE user_credentials SET secret_hash = $1, updated_at = NOW()
	          WHERE user_id = $2 AND is_primary = true`
	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET username = $1, username_updated_at = $2, status = $3, role = $4, updated_at = $5
	          WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, u.Username, u.UsernameUpdatedAt, u.Status, u.Role, u.UpdatedAt, u.ID)
	return err
}
