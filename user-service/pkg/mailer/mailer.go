package mailer

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type Mailer interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type smtpMailer struct {
	host     string
	port     int
	user     string
	password string
	sender   string
	logger   logger.Interface
}

func New(host string, port int, user string, password string, sender string, logger logger.Interface) Mailer {
	if host == "" {
		logger.Info("SMTPHost is empty. Initializing MockMailer (logs emails to console instead of sending)")
		return &mockMailer{logger: logger}
	}
	return &smtpMailer{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		sender:   sender,
		logger:   logger,
	}
}

func (m *smtpMailer) Send(ctx context.Context, to string, subject string, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", m.sender, to, subject, body)
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.user, m.password, m.host)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := smtp.SendMail(addr, auth, m.sender, []string{to}, []byte(msg))
	if err != nil {
		m.logger.Error("failed to send email to %s: %v", to, err)
		return err
	}

	m.logger.Info("successfully sent email to %s (subject: %s)", to, subject)
	return nil
}

type mockMailer struct {
	logger logger.Interface
}

func (m *mockMailer) Send(ctx context.Context, to string, subject string, body string) error {
	m.logger.Info("[MOCK MAILER] Sending Email:\n  To: %s\n  Subject: %s\n  Body: %s", to, subject, body)
	return nil
}
