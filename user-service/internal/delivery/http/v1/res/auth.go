package res

import "user-service/internal/usecase"

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

func ToAuthResponse(res usecase.AuthOutput) AuthResponse {
	return AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		UserID:       res.UserID,
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

type CheckUsernameResponse struct {
	Available bool `json:"available"`
}

func ToCheckUsernameResponse(available bool) CheckUsernameResponse {
	return CheckUsernameResponse{
		Available: available,
	}
}
