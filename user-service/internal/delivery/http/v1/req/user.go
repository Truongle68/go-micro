package req

import (
	"user-service/internal/domain"
	"user-service/internal/usecase"
)

type UpdateProfile struct {
	FullName *string        `json:"full_name"`
	Gender   *domain.Gender `json:"gender"`
	DOB      *string        `json:"dob"`
}

func (req *UpdateProfile) ToUpdatedProfileInput(userID string) usecase.UpdatedProfileInput {
	return usecase.UpdatedProfileInput{
		UserID:   userID,
		FullName: req.FullName,
		Gender:   req.Gender,
		Dob:      req.DOB,
	}
}

type VerifyChangeEmail struct {
	Code string `json:"code" validate:"required"`
}

func (req *VerifyChangeEmail) ToVerifyChangeEmailInput(userID string) usecase.VerifyChangeEmailInput {
	return usecase.VerifyChangeEmailInput{
		UserID: userID,
		Code:   req.Code,
	}
}

type CompleteChangeEmail struct {
	Token    string `json:"token" validate:"required"`
	NewEmail string `json:"new_email" validate:"required,email"`
}

func (req *CompleteChangeEmail) ToCompleteChangeEmailInput(userID string) usecase.CompleteChangeEmailInput {
	return usecase.CompleteChangeEmailInput{
		UserID:   userID,
		Token:    req.Token,
		NewEmail: req.NewEmail,
	}
}

type CompleteChangePhone struct {
	Token    string `json:"token" validate:"required"`
	NewPhone string `json:"new_phone" validate:"required"`
}

func (req *CompleteChangePhone) ToCompleteChangePhoneInput() usecase.CompleteChangePhoneInput {
	return usecase.CompleteChangePhoneInput{
		Token:    req.Token,
		NewPhone: req.NewPhone,
	}
}

type CreateAddress struct {
	Label       string  `json:"label" validate:"required,oneof=home work"`
	AddressLine string  `json:"address_line" validate:"required"`
	Ward        string  `json:"ward"`
	District    string  `json:"district"`
	City        string  `json:"city" validate:"required"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

func (req *CreateAddress) ToCreateAddressInput(userID string) usecase.CreateAddressInput {
	return usecase.CreateAddressInput{
		UserID:      userID,
		Label:       domain.AddressLabel(req.Label),
		AddressLine: req.AddressLine,
		Ward:        req.Ward,
		District:    req.District,
		City:        req.City,
		Lat:         req.Lat,
		Lng:         req.Lng,
	}
}

type UpdateAddress struct {
	ID          string   `json:"id" validate:"required"`
	Label       *string  `json:"label" validate:"omitempty,oneof=home work"`
	AddressLine *string  `json:"address_line"`
	Ward        *string  `json:"ward"`
	District    *string  `json:"district"`
	City        *string  `json:"city"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

func (req *UpdateAddress) ToUpdateAddressInput(userID string) usecase.UpdateAddressInput {
	var label *domain.AddressLabel
	if req.Label != nil {
		l := domain.AddressLabel(*req.Label)
		label = &l
	}
	return usecase.UpdateAddressInput{
		ID:          req.ID,
		UserID:      userID,
		Label:       label,
		AddressLine: req.AddressLine,
		Ward:        req.Ward,
		District:    req.District,
		City:        req.City,
		Lat:         req.Lat,
		Lng:         req.Lng,
	}
}

type SetDefaultAddress struct {
	ID string `json:"id" validate:"required"`
}

type DeleteAddress struct {
	ID string `json:"id" validate:"required"`
}
