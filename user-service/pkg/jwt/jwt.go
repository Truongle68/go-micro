package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"
	"user-service/internal/domain"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type TokenType string

const (
	AccessToken         TokenType = "access"
	RefreshToken        TokenType = "refresh"
	ResetToken          TokenType = "reset"
	VerifyOTPToken      TokenType = "verify_otp"
	ChangeEmailToken    TokenType = "change_email"
	ChangePhoneToken    TokenType = "change_phone"
	ChangePasswordToken TokenType = "change_password"
	EmailLinkToken      TokenType = "email_link"
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

func (c *Claims) JTI() string {
	return c.ID
}

func (c *Claims) ExpiresAtTime() time.Time {
	if c.ExpiresAt == nil {
		return time.Time{}
	}
	return c.ExpiresAt.Time
}

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessTTL  time.Duration
	refreshTTL time.Duration
	actionTTL  time.Duration
}

func New(privateKey, publicKey string, accessTTL, refreshTTL, actionTTL time.Duration) (*JWT, error) {
	prv, err := jwtmanager.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pub, err := jwtmanager.ParsePublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &JWT{
		privateKey: prv,
		publicKey:  pub,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		actionTTL:  actionTTL,
	}, nil
}

var _ TokenService = (*JWT)(nil)

type generateParams struct {
	UserID  string
	Email   string
	Phone   string
	Purpose string
	Role    domain.UserRole
	Type    TokenType
	Expiry  time.Duration
}

func (j *JWT) GenerateAccessToken(userID string, role domain.UserRole) (string, error) {
	return j.generate(generateParams{
		UserID: userID,
		Role:   role,
		Type:   AccessToken,
		Expiry: j.accessTTL,
	})
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	return j.generate(generateParams{
		UserID: userID,
		Type:   RefreshToken,
		Expiry: j.refreshTTL,
	})
}

func (j *JWT) GenerateResetToken(userID string) (string, error) {
	return j.generate(generateParams{
		UserID: userID,
		Type:   ResetToken,
		Expiry: j.actionTTL,
	})
}

func (j *JWT) GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error) {
	return j.generate(generateParams{
		Phone:   phone,
		Purpose: string(purpose),
		Type:    VerifyOTPToken,
		Expiry:  j.actionTTL,
	})
}

func (j *JWT) GenerateChangeEmailToken(userID string) (string, error) {
	return j.generate(generateParams{
		UserID:  userID,
		Purpose: string(domain.VerifyPurposeChangeEmail),
		Type:    ChangeEmailToken,
		Expiry:  j.actionTTL,
	})
}

func (j *JWT) GenerateChangePhoneToken(userID string) (string, error) {
	return j.generate(generateParams{
		UserID:  userID,
		Purpose: string(domain.VerifyPurposeChangePhone),
		Type:    ChangePhoneToken,
		Expiry:  j.actionTTL,
	})
}

func (j *JWT) GenerateChangePasswordToken(userID string) (string, error) {
	return j.generate(generateParams{
		UserID: userID,
		Type:   ChangePasswordToken,
		Expiry: j.actionTTL,
	})
}

func (j *JWT) GenerateEmailLinkToken(userID, email string, purpose domain.EmailLinkPurpose, ttl time.Duration) (string, error) {
	if ttl == 0 {
		ttl = j.actionTTL
	}
	return j.generate(generateParams{
		UserID:  userID,
		Email:   email,
		Purpose: string(purpose),
		Type:    EmailLinkToken,
		Expiry:  ttl,
	})
}

func (j *JWT) generate(params generateParams) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:  params.UserID,
		Email:   params.Email,
		Phone:   params.Phone,
		Role:    string(params.Role),
		Type:    params.Type,
		Purpose: params.Purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(now.Add(params.Expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("jwt - generate - token.SignedString: %w", err)
	}
	return tokenStr, nil
}

func (j *JWT) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, AccessToken)
}

func (j *JWT) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, RefreshToken)
}

func (j *JWT) VerifyResetToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, ResetToken)
}

func (j *JWT) VerifyVerificationToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, VerifyOTPToken)
}

func (j *JWT) VerifyChangeEmailToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, ChangeEmailToken)
}

func (j *JWT) VerifyChangePhoneToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, ChangePhoneToken)
}

func (j *JWT) VerifyChangePasswordToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, ChangePasswordToken)
}

func (j *JWT) VerifyEmailLinkToken(tokenStr string) (*Claims, error) {
	return j.verify(tokenStr, EmailLinkToken)
}

func (j *JWT) verify(tokenStr string, expectedType TokenType) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, ErrUnexpectedSigningMethod
		}
		return j.publicKey, nil
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
