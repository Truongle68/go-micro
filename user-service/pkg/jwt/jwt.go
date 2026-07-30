package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"
	"user-service/internal/domain"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	"github.com/golang-jwt/jwt/v5"
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

func (j *JWT) GenerateAccessToken(userID string, role domain.UserRole) (string, error) {
	return j.generate(userID, "", "", "", role, AccessToken, j.accessTTL)
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	return j.generate(userID, "", "", "", "", RefreshToken, j.refreshTTL)
}

func (j *JWT) GenerateResetToken(userID string) (string, error) {
	return j.generate(userID, "", "", "", "", ResetToken, j.actionTTL)
}

func (j *JWT) GenerateVerificationToken(phone string, purpose domain.VerifyPurpose) (string, error) {
	return j.generate("", "", phone, string(purpose), "", VerifyOTPToken, j.actionTTL)
}

func (j *JWT) GenerateChangeEmailToken(userID string) (string, error) {
	return j.generate(userID, "", "", string(domain.VerifyPurposeChangeEmail), "", ChangeEmailToken, j.actionTTL)
}

func (j *JWT) GenerateChangePhoneToken(userID string) (string, error) {
	return j.generate(userID, "", "", string(domain.VerifyPurposeChangePhone), "", ChangePhoneToken, j.actionTTL)
}

func (j *JWT) GenerateChangePasswordToken(userID string) (string, error) {
	return j.generate(userID, "", "", "", "", ChangePasswordToken, j.actionTTL)
}

func (j *JWT) GenerateEmailLinkToken(userID, email string, purpose domain.EmailLinkPurpose, ttl time.Duration) (string, error) {
	if ttl == 0 {
		ttl = j.actionTTL
	}
	return j.generate(userID, email, "", string(purpose), "", EmailLinkToken, ttl)
}

func (j *JWT) generate(userID, email, phone string, purpose string, role domain.UserRole, tokenType TokenType, expiry time.Duration) (string, error) {
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
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
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
