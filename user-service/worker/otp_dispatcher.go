package worker

import (
	"fmt"
	"user-service/internal/domain"
)

type OTPDispatcher interface {
	DispatchOTP(channel domain.OTPChannel, recipient, code string)
}

type CompositeOTPDispatcher struct {
	email *EmailWorker
	sms   *SMSWorker
}

func NewCompositeOTPDispatcher(email *EmailWorker, sms *SMSWorker) *CompositeOTPDispatcher {
	return &CompositeOTPDispatcher{email: email, sms: sms}
}

var _ OTPDispatcher = (*CompositeOTPDispatcher)(nil)

func (d *CompositeOTPDispatcher) DispatchOTP(channel domain.OTPChannel, recipient, code string) {
	switch channel {
	case domain.OTPChannelEmail:
		d.email.Dispatch(EmailJob{
			Template: "otp_code",
			Subject:  "Your Verification Code",
			To:       recipient,
			Data:     struct{ Code string }{Code: code},
		})
	case domain.OTPChannelSMS:
		d.sms.Dispatch(SMSJob{To: recipient, Body: fmt.Sprintf("Your verification code: %s", code)})
	}
}
