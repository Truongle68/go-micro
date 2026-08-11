package usecase

import (
	"user-service/internal/domain"
	"user-service/pkg/jwt"
)

type RequestOTPInput struct {
	Identifier  string
	Purpose     domain.VerifyPurpose
	ActorUserID string
	ChangeToken string
}

type VerifyOTPInput struct {
	Identifier  string
	Code        string
	Purpose     domain.VerifyPurpose
	ActorUserID string
}

type RegisterInput struct {
	Token    string
	FullName string
	Username string
	Password string
}

type LoginInput struct {
	Identifier    string // can be email, phone or username
	Password      string
	RequiredRoles []domain.UserRole
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

type VerifyPasswordInput struct {
	UserID   string
	Password string
}

type ChangePasswordInput struct {
	UserID            string
	CurrentPassword   string
	NewPassword       string
	ConfirmedPassword string
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type AuthOutput struct {
	AccessToken  jwt.GeneratedTokenOutput
	RefreshToken jwt.GeneratedTokenOutput
	UserID       string
}
