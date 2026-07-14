package usecase

import (
	"context"
	"time"
	"user-service/internal/domain"
)

type UserProfileDTO struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	FullName  string            `json:"full_name"`
	Status    domain.UserStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UpdatedProfileInput struct {
}

type User interface {
	GetProfile(ctx context.Context, id string) (*UserProfileDTO, error)
	UpdateProfile(ctx context.Context, id string, fullName string, phone string) (*UserProfileDTO, error)
}

type RegisterInput struct {
	Token             string
	FullName          string
	Username          string
	Password          string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

type LoginInput struct {
	Identifier string // can be email, phone or username
	Password   string
}

type RequestOTPInput struct {
	Phone   string
	Purpose string
}

type VerifyOTPInput struct {
	Phone   string
	Code    string
	Purpose string
}

type ResetPasswordInput struct {
	Token             string
	NewPassword       string
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type Auth interface {
	RequestOTP(ctx context.Context, in RequestOTPInput) error
	VerifyOTP(ctx context.Context, in VerifyOTPInput) (string, error)
	CompleteRegister(ctx context.Context, in RegisterInput) (AuthOutput, error)
	Login(ctx context.Context, in LoginInput) (AuthOutput, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, in ResetPasswordInput) error
	Logout(ctx context.Context, in LogoutInput) error
}
