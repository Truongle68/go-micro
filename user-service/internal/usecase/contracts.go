package usecase

import (
	"context"
	"time"
	"user-service/internal/domain"
	"user-service/pkg/jwt"
)

type UserProfileDTO struct {
	ID              string            `json:"id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	IsEmailVerified bool              `json:"is_email_verified"`
	Phone           string            `json:"phone"`
	FullName        string            `json:"full_name"`
	Gender          domain.Gender     `json:"gender"`
	DOB             string            `json:"dob"`
	Role            domain.UserRole   `json:"role"`
	Status          domain.UserStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type UpdatedProfileInput struct {
	UserID   string
	Email    *string
	FullName *string
	Gender   *domain.Gender
	Dob      *string
}

type VerifyChangeEmailInput struct {
	UserID string
	Code   string
}

type CompleteChangeEmailInput struct {
	UserID   string
	Token    string
	NewEmail string
}

type CompleteChangePhoneInput struct {
	Token    string
	NewPhone string
}

type AddressDTO struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Label       domain.AddressLabel `json:"label"`
	AddressLine string              `json:"address_line"`
	Ward        string              `json:"ward"`
	District    string              `json:"district"`
	City        string              `json:"city"`
	Lat         float64             `json:"lat"`
	Lng         float64             `json:"lng"`
	IsDefault   bool                `json:"is_default"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type CreateAddressInput struct {
	UserID      string
	Label       domain.AddressLabel
	AddressLine string
	Ward        string
	District    string
	City        string
	Lat         float64
	Lng         float64
}

type UpdateAddressInput struct {
	ID          string
	UserID      string
	Label       *domain.AddressLabel
	AddressLine *string
	Ward        *string
	District    *string
	City        *string
	Lat         *float64
	Lng         *float64
}

type RequestEmailLinkInput struct {
	Email       string                  `json:"email"`
	Purpose     domain.EmailLinkPurpose `json:"purpose"`
	ActorUserID string                  `json:"actor_user_id"`
}

type ConfirmEmailLinkOutput struct {
	Purpose          domain.EmailLinkPurpose `json:"purpose"`
	ChangeEmailToken string                  `json:"change_email_token,omitempty"`
}

type User interface {
	GetProfile(ctx context.Context, id string) (*UserProfileDTO, error)
	UpdateProfile(ctx context.Context, in UpdatedProfileInput) (*UserProfileDTO, error)
	RequestEmailLink(ctx context.Context, in RequestEmailLinkInput) error
	ConfirmEmailLink(ctx context.Context, token string) (*ConfirmEmailLinkOutput, error)

	SendChangeEmailOTP(ctx context.Context, in RequestChangeEmailOTPInput) error
	SendChangePhoneOTP(ctx context.Context, in RequestChangePhoneOTPInput) error
	SendPhoneVerificationOTP(ctx context.Context, userID string) error

	VerifyChangeEmailOTP(ctx context.Context, in VerifyChangeEmailOTPInput) error
	VerifyChangePhoneOTP(ctx context.Context, in VerifyChangePhoneOTPInput) error
	VerifyPhoneVerificationOTP(ctx context.Context, in VerifyPhoneVerificationOTPInput) (string, error)

	GetAddressList(ctx context.Context, userID string) ([]*AddressDTO, error)
	CreateAddress(ctx context.Context, in CreateAddressInput) (*AddressDTO, error)
	UpdateAddress(ctx context.Context, in UpdateAddressInput) (*AddressDTO, error)
	SetDefaultAddress(ctx context.Context, userID string, addressID string) error
	DeleteAddress(ctx context.Context, userID string, addressID string) error
}

type InitOTPInput struct {
	Identifier string
	Channel    domain.OTPChannel
	Purpose    domain.VerifyPurpose
}

type RegisterInput struct {
	Token    string
	FullName string
	Username string
	Password string
}

type AuthOutput struct {
	AccessToken  jwt.GeneratedTokenOutput
	RefreshToken jwt.GeneratedTokenOutput
	UserID       string
}

type LoginInput struct {
	Identifier    string // can be email, phone or username
	Password      string
	RequiredRoles []domain.UserRole
}

type RequestOTPInput struct {
	Identifier  string               `json:"identifier,omitempty"` // phone or email, empty in case verifying the current phone (before typing the new one), getting from credential by userID
	Purpose     domain.VerifyPurpose `json:"purpose"`
	ActorUserID string               `json:"actor_user_id,omitempty"`
	ChangeToken string               `json:"change_token,omitempty"`
}

type RequestChangeEmailOTPInput struct {
	Identifier       string `json:"identifier"`
	ActorUserID      string `json:"actor_user_id"`
	ChangeEmailToken string `json:"change_email_token"`
}

type RequestChangePhoneOTPInput struct {
	Identifier       string `json:"identifier"`
	ActorUserID      string `json:"actor_user_id"`
	ChangePhoneToken string `json:"change_phone_token"`
}

type VerifyOTPInput struct {
	Identifier  string               `json:"identifier,omitempty"`
	Code        string               `json:"code"`
	Purpose     domain.VerifyPurpose `json:"purpose"`
	ActorUserID string               `json:"actor_user_id"`
}

type VerifyChangeEmailOTPInput struct {
	Identifier  string `json:"identifier"`
	Code        string `json:"code"`
	ActorUserID string `json:"actor_user_id"`
}

type VerifyChangePhoneOTPInput struct {
	Identifier  string `json:"identifier"`
	Code        string `json:"code"`
	ActorUserID string `json:"actor_user_id"`
}

type VerifyPhoneVerificationOTPInput struct {
	Code        string `json:"code"`
	ActorUserID string `json:"actor_user_id"`
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

type VerifyPasswordInput struct {
	UserID   string
	Password string
}

type ChangePasswordInput struct {
	UserID            string
	Token             string
	NewPassword       string
	ConfirmedPassword string
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type Auth interface {
	RequestOTP(ctx context.Context, in RequestOTPInput) error
	VerifyOTP(ctx context.Context, in VerifyOTPInput) (token string, exists bool, username string, err error)
	CheckUsernameAvailable(ctx context.Context, username string) (bool, error)
	CompleteRegister(ctx context.Context, in RegisterInput) (AuthOutput, error)
	Login(ctx context.Context, in LoginInput) (AuthOutput, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, in ResetPasswordInput) error
	VerifyPassword(ctx context.Context, in VerifyPasswordInput) (changePasswordToken string, err error)
	ChangePassword(ctx context.Context, in ChangePasswordInput) error
	Logout(ctx context.Context, in LogoutInput) error
	RefreshToken(ctx context.Context, refreshToken string) (AuthOutput, error)
}
