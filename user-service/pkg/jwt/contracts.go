package jwt

import (
	"time"
	"user-service/internal/domain"
)

type GeneratedTokenOutput struct {
	Token string
	Exp   time.Duration
}

type TokenService interface {
	GenerateAccessToken(userID string, role domain.UserRole) (GeneratedTokenOutput, error)
	GenerateRefreshToken(userID string) (GeneratedTokenOutput, error)
	GenerateResetToken(userID string) (string, error)
	GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error)
	GenerateChangeEmailToken(userID string) (string, error)
	GenerateChangePhoneToken(userID string) (string, error)
	GenerateChangePasswordToken(userID string) (string, error)
	GenerateEmailLinkToken(userID, email string, purpose domain.EmailLinkPurpose, ttl time.Duration) (string, error)
	VerifyAccessToken(tokenStr string) (*Claims, error)
	VerifyRefreshToken(tokenStr string) (*Claims, error)
	VerifyResetToken(tokenStr string) (*Claims, error)
	VerifyVerificationToken(tokenStr string) (*Claims, error)
	VerifyChangeEmailToken(tokenStr string) (*Claims, error)
	VerifyChangePhoneToken(tokenStr string) (*Claims, error)
	VerifyChangePasswordToken(tokenStr string) (*Claims, error)
	VerifyEmailLinkToken(tokenStr string) (*Claims, error)
}
