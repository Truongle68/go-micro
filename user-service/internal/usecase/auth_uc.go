package usecase

import (
	"context"
	"fmt"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/jwt"
)

type AuthUC struct {
	userRepo repo.UserRepository
	jwt      *jwt.JWT
}

func NewAuthUC(repo repo.UserRepository, jwt *jwt.JWT) *AuthUC {
	return &AuthUC{
		userRepo: repo,
		jwt:      jwt,
	}
}

var _ Auth = (*AuthUC)(nil)

type RegisterInput struct {
	Username string
	Email    string
	Phone    string
	Password string
	FullName string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

func (uc *AuthUC) Register(ctx context.Context, in RegisterInput) (AuthOutput, error) {
	exists, err := uc.userRepo.ExistsByEmail(ctx, in.Email)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("checking existing email: %w", err)
	}
	if exists {
		return AuthOutput{}, domain.ErrEmailAlreadyExists
	}

	user, err := domain.NewUser(in.Username, in.Email, in.Phone, in.Password, in.FullName)
	if err != nil {
		return AuthOutput{}, err
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return AuthOutput{}, fmt.Errorf("saving user: %w", err)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc *AuthUC) Login(ctx context.Context, in LoginInput) (AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, in.Email)
	if err != nil {
		return AuthOutput{}, domain.ErrInvalidCredentials
	}

	if !user.CheckPassword(in.Password) {
		return AuthOutput{}, domain.ErrInvalidCredentials
	}

	switch user.Status {
	case domain.UserStatusBanned:
		return AuthOutput{}, domain.ErrUserBanned
	case domain.UserStatusInactive:
		return AuthOutput{}, domain.ErrUserInactive
	case domain.UserStatusActive:
		// OK
	default:
		return AuthOutput{}, fmt.Errorf("unknown user status: %s", user.Status)
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthOutput{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}
