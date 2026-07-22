package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID string    `json:"user_id,omitempty"`
	Role   string    `json:"role,omitempty"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

type JWT struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func New(accSecret, refSecret string, accessTTL, refreshTTL time.Duration) *JWT {
	return &JWT{
		accessSecret:  accSecret,
		refreshSecret: refSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

var _ TokenService = (*JWT)(nil)

func (j *JWT) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.accessSecret, AccessToken)
}

func (j *JWT) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return verify(tokenStr, j.refreshSecret, RefreshToken)
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
