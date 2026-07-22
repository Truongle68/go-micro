package domain

import "time"

type VerifyPurpose string

const (
	VerifyPurposeRegister    VerifyPurpose = "register"
	VerifyPurposeChangeEmail VerifyPurpose = "change_email"
	VerifyPurposeVerifyPhone VerifyPurpose = "verify_phone"
	VerifyPurposeChangePhone VerifyPurpose = "change_phone"
)
	
type VerifiedOTPPolicy struct {
	OTPTTL time.Duration
}

var verifiedOTPPolicy = map[VerifyPurpose]VerifiedOTPPolicy{
	VerifyPurposeRegister:    {OTPTTL: 5 * time.Minute},
	VerifyPurposeChangeEmail: {OTPTTL: 5 * time.Minute},
	VerifyPurposeVerifyPhone: {OTPTTL: 5 * time.Minute},
	VerifyPurposeChangePhone: {OTPTTL: 5 * time.Minute},
}

func GetVerifiedOTPPolicy(purpose VerifyPurpose) (VerifiedOTPPolicy, bool) {
	p, ok := verifiedOTPPolicy[purpose]
	return p, ok
}
