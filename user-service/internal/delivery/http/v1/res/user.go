package res

import (
	"time"
	"user-service/internal/usecase"
)

type ProfileResponse struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
	Phone           string    `json:"phone"`
	FullName        string    `json:"full_name"`
	Gender          string    `json:"gender"`
	DOB             string    `json:"dob"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToProfileResponse(res *usecase.UserProfileDTO) ProfileResponse {
	return ProfileResponse{
		ID:              res.ID,
		Username:        res.Username,
		Email:           res.Email,
		IsEmailVerified: res.IsEmailVerified,
		Phone:           res.Phone,
		FullName:        res.FullName,
		Gender:          string(res.Gender),
		DOB:             res.DOB,
		Role:            string(res.Role),
		Status:          string(res.Status),
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}
}

type VerifyPhoneResponse struct {
	ChangePhoneToken string `json:"change_phone_token"`
}

type AddressResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Label       string    `json:"label"`
	AddressLine string    `json:"address_line"`
	Ward        string    `json:"ward"`
	District    string    `json:"district"`
	City        string    `json:"city"`
	FullAddress string    `json:"full_address,omitempty"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToAddressResponse(dto *usecase.AddressDTO) AddressResponse {
	return AddressResponse{
		ID:          dto.ID,
		UserID:      dto.UserID,
		Label:       string(dto.Label),
		AddressLine: dto.AddressLine,
		Ward:        dto.Ward,
		District:    dto.District,
		City:        dto.City,
		FullAddress: dto.FullAddress,
		Lat:         dto.Lat,
		Lng:         dto.Lng,
		IsDefault:   dto.IsDefault,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func ToAddressListResponse(dtos []*usecase.AddressDTO) []AddressResponse {
	list := make([]AddressResponse, len(dtos))
	for i, dto := range dtos {
		list[i] = ToAddressResponse(dto)
	}
	return list
}
