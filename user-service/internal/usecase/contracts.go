package usecase

import (
	"context"
)

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
