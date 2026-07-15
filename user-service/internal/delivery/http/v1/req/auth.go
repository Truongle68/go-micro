package req

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Token             string `json:"token" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,min=8"`
	ConfirmedPassword string `json:"confirmed_password" validate:"required"`
}

type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RequestOTP struct {
	Phone   string `json:"phone" validate:"required"`
	Purpose string `json:"purpose" validate:"required,oneof=register login reset_password"`
}

type VerifyOTP struct {
	Phone   string `json:"phone" validate:"required"`
	OTPCode string `json:"otp_code" validate:"required"`
	Purpose string `json:"purpose" validate:"required,oneof=register login reset_password"`
}

type CompleteRegister struct {
	Token             string `json:"token" validate:"required"`
	FullName          string `json:"full_name" validate:"required"`
	Username          string `json:"username" validate:"required"`
	Password          string `json:"password" validate:"required,min=8"`
	ConfirmedPassword string `json:"confirmed_password" validate:"required"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}