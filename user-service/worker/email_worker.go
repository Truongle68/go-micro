package worker

import (
	"bytes"
	"context"
	"time"
	"user-service/pkg/mailer"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

const _defaultTimeout = 10 * time.Second

type EmailWorker struct {
	mailer  mailer.Mailer
	logger  logger.Interface
	timeout time.Duration
}

type EmailWorkerOption func(*EmailWorker)

func Timeout(t time.Duration) EmailWorkerOption {
	return func(e *EmailWorker) {
		e.timeout = t
	}
}

func NewEmailWorker(mailer mailer.Mailer, logger logger.Interface, opts ...EmailWorkerOption) *EmailWorker {
	w := &EmailWorker{
		mailer:  mailer,
		logger:  logger,
		timeout: _defaultTimeout,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

type EmailJob struct {
	Template string
	Subject  string
	To       string
	Data     any
}

func (w *EmailWorker) Dispatch(job EmailJob) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
		defer cancel()

		defer func() {
			if r := recover(); r != nil {
				w.logger.Error("panic occured in %s email worker: %v", job.Template, r)
			}
		}()

		var buf bytes.Buffer
		if err := templates.ExecuteTemplate(&buf, job.Template+".html", job.Data); err != nil {
			w.logger.Error("failed to render template %s: %v", job.Template, err)
			return
		}

		if err := w.mailer.Send(ctx, job.To, job.Subject, buf.String()); err != nil {
			w.logger.Error("failed to send %s email to %s: %v", job.Template, job.To, err)
			return
		}

		w.logger.Info("send email %s to %s", job.Template, job.To)
	}()
}
