package domain

import "time"

type EmailLinkPurpose string

const (
	EmailLinkPurposeVerifyNew     EmailLinkPurpose = "verify_new"
	EmailLinkPurposeVerifyCurrent EmailLinkPurpose = "verify_current"
)

type EmailLinkPolicy struct {
	TokenTTL time.Duration
}

var emailLinkPolicies = map[EmailLinkPurpose]EmailLinkPolicy{
	EmailLinkPurposeVerifyNew:     {TokenTTL: 15 * time.Minute},
	EmailLinkPurposeVerifyCurrent: {TokenTTL: 10 * time.Minute},
}

func GetEmailLinkPolicy(purpose EmailLinkPurpose) (EmailLinkPolicy, bool) {
	p, ok := emailLinkPolicies[purpose]
	return p, ok
}
