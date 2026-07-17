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
)

type Claims struct {
	UserID  string    `json:"user_id,omitempty"`
	Role    string    `json:"role,omitempty"`
	Phone   string    `json:"phone,omitempty"`
	Type    TokenType `json:"type"`
	Purpose string    `json:"purpose,omitempty"`
	jwt.RegisteredClaims
}

type JWT struct {
	accessSecret    string
	refreshSecret   string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func New(accSecret, refSecret string, accDuration, refDuration time.Duration) *JWT {
	return &JWT{
		accessSecret:    accSecret,
		refreshSecret:   refSecret,
		accessDuration:  accDuration,
		refreshDuration: refDuration,
	}
}

var _ TokenService = (*JWT)(nil)

func (j *JWT) GenerateAccessToken(userID string, role domain.UserRole) (string, error) {
	return generate(userID, "", j.accessSecret, "", role, AccessToken, j.accessDuration)
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	return generate(userID, "", j.refreshSecret, "", "", RefreshToken, j.refreshDuration)
}

func (j *JWT) GenerateResetToken(userID string) (string, error) {
	return generate(userID, "", j.accessSecret, "", "", ResetToken, 15*time.Minute)
}

func (j *JWT) GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error) {
	return generate("", phone, j.accessSecret, purpose, "", VerifyOTPToken, 15*time.Minute)
}

func (j *JWT) GenerateChangeEmailToken(userID string) (string, error) {
	return generate(userID, "", j.accessSecret, "", "", ChangeEmailToken, 15*time.Minute)
}

func (j *JWT) GenerateChangePhoneToken(userID string) (string, error) {
	return generate(userID, "", j.accessSecret, "", "", ChangePhoneToken, 15*time.Minute)
}

func generate(userID, phone, secret string, purpose domain.VerifyPurpose, role domain.UserRole, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:  userID,
		Role:    string(role),
		Phone:   phone,
		Type:    tokenType,
		Purpose: string(purpose),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateToken - token.SignedString: %w", err)
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
