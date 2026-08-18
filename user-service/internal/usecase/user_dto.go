package usecase

import (
	"time"
	"user-service/internal/domain"
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

type RequestEmailLinkInput struct {
	Email       string
	Purpose     domain.EmailLinkPurpose
	ActorUserID string
}

type ConfirmEmailLinkOutput struct {
	Purpose          domain.EmailLinkPurpose
	ChangeEmailToken string
}

type RequestChangeEmailOTPInput struct {
	Identifier       string
	ActorUserID      string
	ChangeEmailToken string
}

type RequestChangePhoneOTPInput struct {
	Identifier       string
	ActorUserID      string
	ChangePhoneToken string
}

type VerifyChangeEmailOTPInput struct {
	Identifier  string
	Code        string
	ActorUserID string
}

type VerifyChangePhoneOTPInput struct {
	Identifier  string
	Code        string
	ActorUserID string
}

type VerifyPhoneVerificationOTPInput struct {
	Code        string
	ActorUserID string
}

type AddressDTO struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Label       domain.AddressLabel `json:"label"`
	AddressLine string              `json:"address_line"`
	Ward        string              `json:"ward"`
	District    string              `json:"district"`
	City        string              `json:"city"`
	FullAddress string              `json:"full_address,omitempty"`
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
