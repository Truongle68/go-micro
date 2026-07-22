package jwt

import (
	"errors"
	"fmt"
	"time"
	"user-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type TokenType string

const (
	AccessToken      TokenType = "access"
	RefreshToken     TokenType = "refresh"
	ResetToken       TokenType = "reset"
	VerifyOTPToken   TokenType = "verify_otp"
	ChangeEmailToken TokenType = "change_email"
	ChangePhoneToken TokenType = "change_phone"
	EmailLinkToken   TokenType = "email_link"
)

type Claims struct {
	UserID  string    `json:"user_id,omitempty"`
	Role    string    `json:"role,omitempty"`
	Phone   string    `json:"phone,omitempty"`
	Email   string    `json:"email,omitempty"`
	Type    TokenType `json:"type"`
	Purpose string    `json:"purpose,omitempty"`
	jwt.RegisteredClaims
}

type JWT struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	actionTTL     time.Duration
}

func New(accSecret, refSecret string, accessTTL, refreshTTL, actionTTL time.Duration) *JWT {
	return &JWT{
		accessSecret:  accSecret,
		refreshSecret: refSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		actionTTL:     actionTTL,
	}
}

var _ TokenService = (*JWT)(nil)

func (j *JWT) GenerateAccessToken(userID string, role domain.UserRole) (string, error) {
	return generate(userID, "", "", j.accessSecret, "", role, AccessToken, j.accessTTL)
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	return generate(userID, "", "", j.refreshSecret, "", "", RefreshToken, j.refreshTTL)
}

func (j *JWT) GenerateResetToken(userID string) (string, error) {
	return generate(userID, "", "", j.accessSecret, "", "", ResetToken, j.actionTTL)
}

func (j *JWT) GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error) {
	return generate("", "", phone, j.accessSecret, string(purpose), "", VerifyOTPToken, j.actionTTL)
}

func (j *JWT) GenerateChangeEmailToken(userID string) (string, error) {
	return generate(userID, "", "", j.accessSecret, string(domain.VerifyPurposeChangeEmail), "", ChangeEmailToken, j.actionTTL)
}

func (j *JWT) GenerateChangePhoneToken(userID string) (string, error) {
	return generate(userID, "", "", j.accessSecret, string(domain.VerifyPurposeChangePhone), "", ChangePhoneToken, j.actionTTL)
}

func (j *JWT) GenerateEmailLinkToken(userID, email string, purpose domain.EmailLinkPurpose, ttl time.Duration) (string, error) {
	if ttl == 0 {
		ttl = j.actionTTL
	}
	return generate(userID, email, "", j.accessSecret, string(purpose), "", EmailLinkToken, ttl)
}

func generate(userID, email, phone, secret string, purpose string, role domain.UserRole, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:  userID,
		Email:   email,
		Phone:   phone,
		Role:    string(role),
		Type:    tokenType,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt - generate - token.SignedString: %w", err)
	}
	return tokenStr, nil
}

func (j *JWT) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, AccessToken)
}

func (j *JWT) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.refreshSecret, RefreshToken)
}

func (j *JWT) VerifyResetToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, ResetToken)
}

func (j *JWT) VerifyVerificationToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, VerifyOTPToken)
}

func (j *JWT) VerifyChangeEmailToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, ChangeEmailToken)
}

func (j *JWT) VerifyChangePhoneToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, ChangePhoneToken)
}

func (j *JWT) VerifyEmailLinkToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, EmailLinkToken)
}

func verify(tokenStr, secret string, expectedType TokenType) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt - verify - jwt.ParseWithClaims: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("jwt - verify - expected %q token, got %q", expectedType, claims.Type)
	}

	return claims, nil
}
