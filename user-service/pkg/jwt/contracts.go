package jwt

import (
	"time"
	"user-service/internal/domain"
)

type TokenService interface {
	GenerateAccessToken(userID string, role domain.UserRole) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateResetToken(userID string) (string, error)
	GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error)
	GenerateChangeEmailToken(userID string) (string, error)
	GenerateChangePhoneToken(userID string) (string, error)
	GenerateEmailLinkToken(userID, email string, purpose domain.EmailLinkPurpose, ttl time.Duration) (string, error)
	VerifyAccessToken(tokenStr string) (*Claims, error)
	VerifyRefreshToken(tokenStr string) (*Claims, error)
	VerifyResetToken(tokenStr string) (*Claims, error)
	VerifyVerificationToken(tokenStr string) (*Claims, error)
	VerifyChangeEmailToken(tokenStr string) (*Claims, error)
	VerifyChangePhoneToken(tokenStr string) (*Claims, error)
	VerifyEmailLinkToken(tokenStr string) (*Claims, error)
}
