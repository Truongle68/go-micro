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

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

var PortalRole []UserRole = []UserRole{
	UserRoleAdmin,
}

type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	UsernameUpdatedAt time.Time  `json:"username_updated_at"`
	Status            UserStatus `json:"status"`
	Role              UserRole   `json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func isMatch[T comparable](left, right T) bool {
	return left == right
}

func IsConfirmMatch(password, confirm string) bool {
	return isMatch(password, confirm)
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
		Role:              UserRoleUser,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u, nil
}
