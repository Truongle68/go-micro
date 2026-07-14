package req

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

type UpdateProfile struct {
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}
