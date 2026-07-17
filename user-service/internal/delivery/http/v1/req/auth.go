package req

import (
	"user-service/internal/domain"
	"user-service/internal/usecase"
)

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Token             string `json:"token" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,min=8"`
	ConfirmedPassword string `json:"confirmed_password" validate:"required"`
}

func (req *ResetPassword) ToResetPasswordInput() usecase.ResetPasswordInput {
	return usecase.ResetPasswordInput{
		Token: req.Token,
		NewPassword: req.NewPassword,
	}
}

type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (req *Logout) ToLogoutInput(accessToken string) usecase.LogoutInput {
	return usecase.LogoutInput{}
}

type RequestOTP struct {
	Phone string `json:"phone" validate:"required"`
}

func (req *RequestOTP) ToRequestOTPInput(purpose domain.VerifyPurpose) usecase.RequestOTPInput {
	return usecase.RequestOTPInput{
		Phone:   req.Phone,
		Purpose: purpose,
	}
}

type VerifyOTP struct {
	Phone   string `json:"phone" validate:"required"`
	OTPCode string `json:"otp_code" validate:"required"`
}

func (req *VerifyOTP) ToVerifyOTPInput(purpose domain.VerifyPurpose) usecase.VerifyOTPInput {
	return usecase.VerifyOTPInput{
		Phone:   req.Phone,
		Code:    req.OTPCode,
		Purpose: purpose,
	}
}

type CompleteRegister struct {
	Token             string `json:"token" validate:"required"`
	FullName          string `json:"full_name" validate:"required"`
	Username          string `json:"username" validate:"required"`
	Password          string `json:"password" validate:"required,min=8"`
	ConfirmedPassword string `json:"confirmed_password" validate:"required"`
}

func (req *CompleteRegister) ToRegisterInput() usecase.RegisterInput {
	return usecase.RegisterInput{
		Token:    req.Token,
		FullName: req.FullName,
		Username: req.Username,
		Password: req.Password,
	}
}

type Login struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}

func (req *Login) ToLoginInput(requiredRoles []domain.UserRole) usecase.LoginInput {
	return usecase.LoginInput{
		Identifier:    req.Identifier,
		Password:      req.Password,
		RequiredRoles: requiredRoles,
	}
}
