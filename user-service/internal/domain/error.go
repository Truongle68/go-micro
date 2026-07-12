package domain

import "errors"

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserBanned         = errors.New("user account is banned")
	ErrUserInactive       = errors.New("user account is inactive")
)
