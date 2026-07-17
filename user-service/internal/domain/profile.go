package domain

import "time"

type Gender string

const (
	GenderUnknown Gender = "unknown"
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
)

type Profile struct {
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	Gender    Gender    `json:"gender"`
	DOB       string    `json:"dob"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileParams struct {
	FullName *string
	Gender   *Gender
	Dob      *string
}

func NewProfile(userID, fullName string) *Profile {
	return &Profile{
		UserID:    userID,
		FullName:  fullName,
		UpdatedAt: time.Now(),
	}
}

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther:
		return true
	}
	return false
}

func (p *Profile) ApplyUpdate(params UpdateProfileParams) error {
	change := false
	if params.FullName != nil {
		if *params.FullName == "" {
			return ErrInvalidFullName
		}
		p.FullName = *params.FullName
		change = true
	}

	if params.Gender != nil {
		if !params.Gender.IsValid() {
			return ErrInvalidGender
		}
		p.Gender = *params.Gender
		change = true
	}
	if params.Dob != nil {
		if _, err := time.Parse("2026-01-02", *params.Dob); err != nil {
			return ErrInvalidDob
		}
		p.DOB = *params.Dob
		change = true
	}
	if !change {
		return ErrNoFieldsToUpdate
	}
	p.UpdatedAt = time.Now()
	return nil
}
