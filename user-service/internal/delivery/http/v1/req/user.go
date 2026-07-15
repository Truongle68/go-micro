package req

type UpdateProfile struct {
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}
