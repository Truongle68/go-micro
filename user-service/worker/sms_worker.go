// worker/sms_worker.go
package worker

import (
	"context"
	"time"
	"user-service/pkg/sms" // your SMS gateway client, mirroring pkg/mailer

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type SMSWorker struct {
	client  sms.Client
	logger  logger.Interface
	timeout time.Duration
}

func NewSMSWorker(client sms.Client, logger logger.Interface) *SMSWorker {
	return &SMSWorker{client: client, logger: logger, timeout: _defaultTimeout}
}

type SMSJob struct { 
	To   string
	Body string
}

func (w *SMSWorker) Dispatch(job SMSJob) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				w.logger.Error("panic occurred in SMS worker: %v", r)
			}
		}()

		if err := w.client.Send(ctx, job.To, job.Body); err != nil {
			w.logger.Error("failed to send SMS to %s: %v", job.To, err)
			return
		}
		w.logger.Info("sent SMS to %s", job.To)
	}()
}