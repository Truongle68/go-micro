package req

import (
	"user-service/internal/domain"
	"user-service/internal/usecase"
)

type UpdateProfile struct {
	Email    *string        `json:"email"`
	FullName *string        `json:"full_name"`
	Gender   *domain.Gender `json:"gender"`
	DOB      *string        `json:"dob"`
}

func (req *UpdateProfile) ToInput(userID string) usecase.UpdatedProfileInput {
	return usecase.UpdatedProfileInput{
		UserID:   userID,
		Email:    req.Email,
		FullName: req.FullName,
		Gender:   req.Gender,
		Dob:      req.DOB,
	}
}

type RequestEmailLink struct {
	Email   string                  `json:"email" validate:"required,email"`
	Purpose domain.EmailLinkPurpose `json:"purpose"`
}

func (req *RequestEmailLink) ToInput(actorUserID string) usecase.RequestEmailLinkInput {
	p := req.Purpose
	if p == "" {
		p = domain.EmailLinkPurposeVerifyNew
	}
	return usecase.RequestEmailLinkInput{
		Email:       req.Email,
		Purpose:     p,
		ActorUserID: actorUserID,
	}
}

type ChangeEmail struct {
	Identifier       string `json:"identifier" validate:"required,email"`
	ChangeEmailToken string `json:"change_email_token" validate:"required"`
}

func (req *ChangeEmail) ToInput(actorUserID string) usecase.RequestChangeEmailOTPInput {
	return usecase.RequestChangeEmailOTPInput{
		Identifier:       req.Identifier,
		ActorUserID:      actorUserID,
		ChangeEmailToken: req.ChangeEmailToken,
	}
}

type ChangeEmailConfirm struct {
	Identifier string `json:"identifier" validate:"required,email"`
	Code       string `json:"code" validate:"required"`
}

func (req *ChangeEmailConfirm) ToInput(actorUserID string) usecase.VerifyChangeEmailOTPInput {
	return usecase.VerifyChangeEmailOTPInput{
		Identifier:  req.Identifier,
		Code:        req.Code,
		ActorUserID: actorUserID,
	}
}

type VerifyPhone struct {
	Code string `json:"code" validate:"required"`
}

func (req *VerifyPhone) ToInput(actorUserID string) usecase.VerifyPhoneVerificationOTPInput {
	return usecase.VerifyPhoneVerificationOTPInput{
		Code:        req.Code,
		ActorUserID: actorUserID,
	}
}

type ChangePhone struct {
	Phone            string `json:"phone" validate:"required"`
	ChangePhoneToken string `json:"change_phone_token" validate:"required"`
}

func (req *ChangePhone) ToInput(actorUserID string) usecase.RequestChangePhoneOTPInput {
	return usecase.RequestChangePhoneOTPInput{
		Identifier:       req.Phone,
		ActorUserID:      actorUserID,
		ChangePhoneToken: req.ChangePhoneToken,
	}
}

type ChangePhoneVerify struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (req *ChangePhoneVerify) ToInput(actorUserID string) usecase.VerifyChangePhoneOTPInput {
	return usecase.VerifyChangePhoneOTPInput{
		Identifier:  req.Phone,
		Code:        req.Code,
		ActorUserID: actorUserID,
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

func (req *CreateAddress) ToInput(userID string) usecase.CreateAddressInput {
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
	Label       *string  `json:"label" validate:"omitempty,oneof=home work"`
	AddressLine *string  `json:"address_line"`
	Ward        *string  `json:"ward"`
	District    *string  `json:"district"`
	City        *string  `json:"city"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

func (req *UpdateAddress) ToInput(id, userID string) usecase.UpdateAddressInput {
	var label *domain.AddressLabel
	if req.Label != nil {
		l := domain.AddressLabel(*req.Label)
		label = &l
	}
	return usecase.UpdateAddressInput{
		ID:          id,
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
