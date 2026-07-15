package domain

import "errors"

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrNotMatchPassword   = errors.New("confirmed password is not match")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("user account is banned")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidOTP         = errors.New("invalid OTP code")
	ErrOTPExpired         = errors.New("OTP code expired")
	ErrEmptyUsername      = errors.New("username cannot be empty")
	ErrEmptyEmail         = errors.New("email cannot be empty")
	ErrEmptyPhone         = errors.New("phone cannot be empty")
	ErrUsernameExists     = errors.New("username already exists")
	ErrPhoneAlreadyExists = errors.New("phone number already registered")
)
