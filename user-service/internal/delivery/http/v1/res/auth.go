package res

import (
	"time"
	"user-service/internal/usecase"
)

type AuthResponse struct {
	UserID     string        `json:"user_id"`
	AccessExp  time.Duration `json:"access_exp"`
	RefreshExp time.Duration `json:"refresh_exp"`
}

func ToAuthResponse(res usecase.AuthOutput) AuthResponse {
	return AuthResponse{
		UserID:     res.UserID,
		AccessExp:  res.AccessToken.Exp,
		RefreshExp: res.RefreshToken.Exp,
	}
}

type VerifyOTPResponse struct {
	VerifyOTPToken string `json:"verify_otp_token"`
	UserExists     bool   `json:"user_exists"`
	Username       string `json:"username,omitempty"`
}

func ToVerifyOTPResponse(token string, exists bool, username string) VerifyOTPResponse {
	return VerifyOTPResponse{
		VerifyOTPToken: token,
		UserExists:     exists,
		Username:       username,
	}
}

type ForgotPasswordResponse struct {
	ResetToken string `json:"reset_token"`
}

func ToForgotPasswordResponse(token string) ForgotPasswordResponse {
	return ForgotPasswordResponse{
		ResetToken: token,
	}
}

type VerifyPasswordResponse struct {
	ChangePasswordToken string `json:"change_password_token"`
}

func ToVerifyPasswordResponse(token string) VerifyPasswordResponse {
	return VerifyPasswordResponse{
		ChangePasswordToken: token,
	}
}

type CheckUsernameResponse struct {
	Available bool `json:"available"`
}

func ToCheckUsernameResponse(available bool) CheckUsernameResponse {
	return CheckUsernameResponse{
		Available: available,
	}
}
