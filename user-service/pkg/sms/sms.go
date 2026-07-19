package sms

import (
	"context"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type Client interface {
	Send(ctx context.Context, to string, body string) error
}

type mockSMS struct {
	logger logger.Interface
}

func NewMockSMS(logger logger.Interface) Client {
	return &mockSMS{logger: logger}
}

func (s *mockSMS) Send(ctx context.Context, to string, body string) error {
	s.logger.Info("[MOCK SMS] Sending SMS:\n  To: %s\n  Body: %s", to, body)
	return nil
}
