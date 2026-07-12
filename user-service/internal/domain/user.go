package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewUser(username, email, phone, plainPassword, fullName string) (*User, error) {
	if email == "" {
		return nil, ErrEmailRequired
	}
	if len(plainPassword) < 8 {
		return nil, ErrWeakPassword
	}

	hash, err := hashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Username:     username,
		Email:        email,
		Phone:        phone,
		PasswordHash: hash,
		FullName:     fullName,
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) CheckPassword(pass string) bool {
	return checkPassword(pass, u.PasswordHash)
}
