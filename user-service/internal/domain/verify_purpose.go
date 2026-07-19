package domain

type VerifyPurpose string

const (
	VerifyPurposeRegister    VerifyPurpose = "register"
	VerifyPurposeChangeEmail VerifyPurpose = "change_email"
)
