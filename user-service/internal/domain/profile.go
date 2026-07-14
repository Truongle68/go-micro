package domain

import "time"

type Profile struct {
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	Gender    string    `json:"gender"`
	DOB       string    `json:"dob"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProfile(userID, fullName string) *Profile {
	return &Profile{
		UserID:    userID,
		FullName:  fullName,
		UpdatedAt: time.Now(),
	}
}
