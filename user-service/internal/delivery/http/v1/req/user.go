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
