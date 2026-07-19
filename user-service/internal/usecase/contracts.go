package usecase

import (
	"context"
	"time"
	"user-service/internal/domain"
)

type UserProfileDTO struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	FullName  string            `json:"full_name"`
	Gender    domain.Gender     `json:"gender"`
	DOB       string            `json:"dob"`
	Role      domain.UserRole   `json:"role"`
	Status    domain.UserStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UpdatedProfileInput struct {
	UserID   string
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

type User interface {
	GetProfile(ctx context.Context, id string) (*UserProfileDTO, error)
	UpdateProfile(ctx context.Context, in UpdatedProfileInput) (*UserProfileDTO, error)
	RequestChangeEmailOTP(ctx context.Context, userID string) error
	VerifyChangeEmailOTP(ctx context.Context, in VerifyChangeEmailInput) (token string, err error)
	CompleteChangeEmail(ctx context.Context, in CompleteChangeEmailInput) error
	RequestChangePhoneLink(ctx context.Context, userID string) error
	CompleteChangePhone(ctx context.Context, in CompleteChangePhoneInput) error
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
	AccessToken  string
	RefreshToken string
	UserID       string
}

type LoginInput struct {
	Identifier    string // can be email, phone or username
	Password      string
	RequiredRoles []domain.UserRole
}

type RequestOTPInput struct {
	Phone   string
	Purpose domain.VerifyPurpose
}

type VerifyOTPInput struct {
	Phone   string
	Code    string
	Purpose domain.VerifyPurpose
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
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
	Logout(ctx context.Context, in LogoutInput) error
}
