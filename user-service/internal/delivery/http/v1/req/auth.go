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
