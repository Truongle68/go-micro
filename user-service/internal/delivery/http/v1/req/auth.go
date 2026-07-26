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

func (req *ResetPassword) ToInput() usecase.ResetPasswordInput {
	return usecase.ResetPasswordInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}
}

type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (req *Logout) ToInput(accessToken string) usecase.LogoutInput {
	return usecase.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RequestOTP struct {
	Identifier string `json:"identifier"`
	Phone      string `json:"phone"`
}

func (req *RequestOTP) ToInput(purpose domain.VerifyPurpose) usecase.RequestOTPInput {
	id := req.Identifier
	if id == "" {
		id = req.Phone
	}
	return usecase.RequestOTPInput{
		Identifier: id,
		Purpose:    purpose,
	}
}

type VerifyOTP struct {
	Identifier string `json:"identifier"`
	Phone      string `json:"phone"`
	OTPCode    string `json:"otp_code" validate:"required"`
}

func (req *VerifyOTP) ToInput(purpose domain.VerifyPurpose) usecase.VerifyOTPInput {
	id := req.Identifier
	if id == "" {
		id = req.Phone
	}
	return usecase.VerifyOTPInput{
		Identifier: id,
		Code:       req.OTPCode,
		Purpose:    purpose,
	}
}

type CompleteRegister struct {
	Token             string `json:"token" validate:"required"`
	FullName          string `json:"full_name" validate:"required"`
	Username          string `json:"username" validate:"required"`
	Password          string `json:"password" validate:"required,min=8"`
	ConfirmedPassword string `json:"confirmed_password" validate:"required"`
}

func (req *CompleteRegister) ToInput() usecase.RegisterInput {
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

func (req *Login) ToInput(requiredRoles []domain.UserRole) usecase.LoginInput {
	return usecase.LoginInput{
		Identifier:    req.Identifier,
		Password:      req.Password,
		RequiredRoles: requiredRoles,
	}
}
