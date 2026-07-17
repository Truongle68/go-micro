package res

import (
	"time"
	"user-service/internal/usecase"
)

type ProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FullName  string    `json:"full_name"`
	Gender    string    `json:"gender"`
	DOB       string    `json:"dob"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToProfileResponse(res *usecase.UserProfileDTO) ProfileResponse {
	return ProfileResponse{
		ID:        res.ID,
		Username:  res.Username,
		Email:     res.Email,
		Phone:     res.Phone,
		FullName:  res.FullName,
		Gender:    string(res.Gender),
		DOB:       res.DOB,
		Role:      string(res.Role),
		Status:    string(res.Status),
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}

type VerifyChangePhoneResponse struct {
	Token string `json:"token"`
}
