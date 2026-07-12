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
	ResetToken   TokenType = "reset"
)

type Claims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`
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

func (j *JWT) GenerateAccessToken(userID string) (string, error) {
	return generate(userID, j.accessSecret, AccessToken, j.accessDuration)
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	return generate(userID, j.refreshSecret, RefreshToken, j.refreshDuration)
}

func (j *JWT) GenerateResetToken(userID string) (string, error) {
	return generate(userID, j.accessSecret, ResetToken, 15*time.Minute)
}

func generate(userID, secret string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Type:   tokenType,
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
		return nil, fmt.Errorf("jwt - verify - expected %q token, got %q: %w", expectedType, claims.Type, err)
	}

	return claims, err
}
