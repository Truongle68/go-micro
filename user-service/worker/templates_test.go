package worker

import (
	"bytes"
	"testing"
)

func TestTemplatesRendering(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     any
	}{
		{
			name:     "change_email_otp",
			template: "change_email_otp.html",
			data:     struct{ Code string }{Code: "123456"},
		},
		{
			name:     "otp_code",
			template: "otp_code.html",
			data:     struct{ Code string }{Code: "654321"},
		},
		{
			name:     "verify_email",
			template: "verify_email.html",
			data:     struct{ Link string }{Link: "http://example.com"},
		},
		{
			name:     "forgot_password",
			template: "forgot_password.html",
			data:     struct{ ResetLink string }{ResetLink: "http://example.com/reset"},
		},
		{
			name:     "change_phone",
			template: "change_phone.html",
			data:     struct{ ResetLink string }{ResetLink: "http://example.com/phone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := templates.ExecuteTemplate(&buf, tt.template, tt.data)
			if err != nil {
				t.Fatalf("failed to execute template %s: %v", tt.template, err)
			}
			if buf.Len() == 0 {
				t.Fatalf("rendered template %s is empty", tt.template)
			}
		})
	}
}
