package domain

import "time"

type OTPSession struct {
	Seed                string
	Phone               string
	Purpose             string
	NeedCaptcha         bool
	ResendCooldownUntil time.Time
	CreatedAt           time.Time
}
