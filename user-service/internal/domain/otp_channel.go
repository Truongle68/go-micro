package domain

type OTPChannel string

const (
	OTPChannelSMS   OTPChannel = "sms"
	OTPChannelEmail OTPChannel = "email"
)
