package domain

import (
	"time"
)

type UserStatus string

const (
	UserStatusVerified   UserStatus = "verified"
	UserStatusUnverified UserStatus = "unverified"
	UserStatusBanned     UserStatus = "banned"
)

type CredentialType string

const (
	CredentialTypePhone  CredentialType = "phone"
	CredentialTypeEmail  CredentialType = "email"
	CredentialTypeGoogle CredentialType = "google"
)

type VerifyPurpose string

const (
	VerifyPurposeRegister      VerifyPurpose = "register"
	VerifyPurposeLogin         VerifyPurpose = "login"
	VerifyPurposeResetPassword VerifyPurpose = "reset_password"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	UsernameUpdatedAt time.Time  `json:"username_updated_at"`
	Status            UserStatus `json:"status"`
	Role              UserRole   `json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type UserCredential struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	Type       CredentialType `json:"type"`
	Identifier string         `json:"identifier"`
	SecretHash string         `json:"-"`
	IsVerified bool           `json:"is_verified"`
	IsPrimary  bool           `json:"is_primary"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func ValidatePassword(pass string) error {
	if len(pass) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func isMatch[T comparable](left, right T) bool {
	return left == right
}

func IsConfirmMatch(password, confirm string) bool {
	return isMatch(password, confirm)
}

func (c *UserCredential) CheckPassword(pass string) bool {
	return checkPassword(pass, c.SecretHash)
}

func NewUser(username string, opts ...UserOption) (*User, error) {
	if username == "" {
		return nil, ErrEmptyUsername
	}
	now := time.Now()
	u := &User{
		Username:          username,
		UsernameUpdatedAt: now,
		Status:            UserStatusVerified,
		Role:              UserRoleCustomer,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u, nil
}

func NewUserCredential(userID string, credType CredentialType, identifier, plainPassword string, isVerified, isPrimary bool) (*UserCredential, error) {
	var hash string
	if plainPassword != "" {
		if err := ValidatePassword(plainPassword); err != nil {
			return nil, err
		}
		var err error
		hash, err = HashPassword(plainPassword)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	return &UserCredential{
		UserID:     userID,
		Type:       credType,
		Identifier: identifier,
		SecretHash: hash,
		IsVerified: isVerified,
		IsPrimary:  isPrimary,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
