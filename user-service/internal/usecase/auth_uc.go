package usecase

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pkg/jwt"

	"github.com/redis/go-redis/v9"
)

type AuthUC struct {
	userRepo    repo.UserRepository
	jwt         *jwt.JWT
	redisClient *redis.Client
}

func NewAuthUC(repo repo.UserRepository, jwt *jwt.JWT, redisClient *redis.Client) *AuthUC {
	return &AuthUC{
		userRepo:    repo,
		jwt:         jwt,
		redisClient: redisClient,
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

func (uc *AuthUC) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	resetToken, err := uc.jwt.GenerateResetToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	return resetToken, nil
}

func (uc *AuthUC) ResetPassword(ctx context.Context, token string, newPassword string) error {
	claims, err := uc.jwt.VerifyResetToken(token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if len(newPassword) < 8 {
		return domain.ErrWeakPassword
	}

	hash, err := domain.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = uc.userRepo.UpdatePassword(ctx, claims.UserID, hash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// invalidate user profile cache on password reset
	cacheKey := fmt.Sprintf("user:profile:%s", claims.UserID)
	uc.redisClient.Del(ctx, cacheKey)

	return nil
}

func (uc *AuthUC) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	// blacklist access token
	accClaims, err := uc.jwt.VerifyAccessToken(accessToken)
	if err == nil {
		accTtl := time.Until(accClaims.ExpiresAt.Time)
		if accTtl > 0 {
			err = uc.redisClient.Set(ctx, "blacklist:"+accessToken, "1", accTtl).Err()
			if err != nil {
				return fmt.Errorf("failed to blacklist access token: %w", err)
			}
		}
	}

	// blacklist refresh token
	refClaims, err := uc.jwt.VerifyRefreshToken(refreshToken)
	if err == nil {
		refTtl := time.Until(refClaims.ExpiresAt.Time)
		if refTtl > 0 {
			err = uc.redisClient.Set(ctx, "blacklist:"+refreshToken, "1", refTtl).Err()
			if err != nil {
				return fmt.Errorf("failed to blacklist refresh token: %w", err)
			}
		}
	}

	return nil
}
