package usecase

import (
	"context"
	"user-service/internal/domain"
)

type UserUsecase interface {
	GetProfile(ctx context.Context, id string) (*domain.User, error)
	UpdateProfile(ctx context.Context, id string, fullName string, phone string) (*domain.User, error)
}

type Auth interface {
	Register(ctx context.Context, in RegisterInput) (AuthOutput, error)
	Login(ctx context.Context, in LoginInput) (AuthOutput, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
	Logout(ctx context.Context, accessToken string, refreshToken string) error
}
