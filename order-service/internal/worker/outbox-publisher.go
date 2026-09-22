package worker

import (
	"context"
	"fmt"
	"order-service/internal/domain"
	"sync"
	"time"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

const (
	_defaultPollInterval = 5 * time.Second
	_defaultBatchSize    = 50
	_defaultMaxRetries   = 5
)

// OutboxStore defines the persistent capabilities required by the worker.
type OutboxStore interface {
	GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkPublished(ctx context.Context, ids []string) error
	IncrementRetry(ctx context.Context, id string) error
}

// Transactor defines transaction execution.
type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// EventPublisher defines the message broker publishing interface.
type EventPublisher interface {
	Publish(ctx context.Context, event rabbitmq.Event) error
}

// Option configures OutboxPublisherWorker.
type Option func(*OutboxPublisherWorker)

// PollInterval sets how often the worker polls for unpublished event.
func PollInterval(d time.Duration) Option {
	return func(w *OutboxPublisherWorker) {
		if d > 0 {
			w.pollInterval = d
		}
	}
}

// BatchSize sets maximum number of events processed per cycle.
func BatchSize(s int) Option {
	return func(w *OutboxPublisherWorker) {
		if s > 0 {
			w.batchSize = s
		}
	}
}

// MaxRetries sets maximum retries before an event is skipped or flagged.
func MaxRetries(retries int) Option {
	return func(w *OutboxPublisherWorker) {
		if retries > 0 {
			w.maxRetries = retries
		}
	}
}

// OutboxPublisherWorker polls unpublished outbox records and publishes them to RabbitMQ.
type OutboxPublisherWorker struct {
	repo         OutboxStore
	transactor   Transactor
	publisher    EventPublisher
	logger       logger.Interface
	pollInterval time.Duration
	batchSize    int
	maxRetries   int

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewOutBoxPublisherWorker creates a new OutboxPublisherWorker instance.
func NewOutBoxPublisherWorker(
	repo OutboxStore,
	transactor Transactor,
	publisher EventPublisher,
	l logger.Interface,
	opts ...Option,
) *OutboxPublisherWorker {
	w := &OutboxPublisherWorker{
		repo:         repo,
		transactor:   transactor,
		publisher:    publisher,
		logger:       l,
		pollInterval: _defaultPollInterval,
		batchSize:    _defaultBatchSize,
		maxRetries:   _defaultMaxRetries,
		stopCh:       make(chan struct{}),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Start begins the background polling loop. It blocks until ctx is done or Stop() is called.
func (w *OutboxPublisherWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	w.logger.Info("OutboxPublisherWorker started (poll_interval: %v, batch_size: %d)", w.pollInterval, w.batchSize)

	// Process immediately once on start.
	if err := w.processBatchSafe(ctx); err != nil {
		w.logger.Error("OutboxPublisherWorker initial cycle error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("OutboxPublisherWorker context cancelled, stopping...")
			return
		case <-w.stopCh:
			w.logger.Info("OutboxPublisherWorker stop signal received, stopping...")
			return
		case <-ticker.C:
			if err := w.processBatchSafe(ctx); err != nil {
				w.logger.Error("OutboxPublisherWorker initial cycle error: %v", err)
			}
		}
	}
}

// Stop signals the worker to stop and wait for in-flight processing to complete.
func (w *OutboxPublisherWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *OutboxPublisherWorker) processBatchSafe(ctx context.Context) error {
	count, err := w.ProcessBatch(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		w.logger.Info("OutboxPublisherWorker processed %d events successfully", count)
	}
	return nil
}

// ProcessBatch fetches and publishes one batch of unpublished events within a transaction.
func (w *OutboxPublisherWorker) ProcessBatch(ctx context.Context) (int, error) {
	var processedCount int

	err := w.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		events, err := w.repo.GetUnpublished(txCtx, w.batchSize)
		if err != nil {
			return fmt.Errorf("fetching unpublished outbox events: %w", err)
		}

		if len(events) == 0 {
			return nil
		}

		publishedIDs := make([]string, 0, len(events))
		for _, evt := range events {
			if evt.RetryCount >= w.maxRetries {
				w.logger.Info("OutboxPublisherWorker event %s exceeded max retries (%d), skipping", evt.ID, w.maxRetries)
				continue
			}

			pubEvent := rabbitmq.Event{
				Type:      evt.EventType,
				Payload:   evt.Payload,
				Timestamp: evt.CreatedAt,
			}

			if err := w.publisher.Publish(txCtx, pubEvent); err != nil {
				w.logger.Error("OutboxPublisherWorker failed to publish event %s: %v", evt.ID, err)
				if retryErr := w.repo.IncrementRetry(txCtx, evt.ID); retryErr != nil {
					w.logger.Error("OutboxPublisherWorker failed to increment retry for %s: %v", evt.ID, retryErr)
				}
				continue
			}

			publishedIDs = append(publishedIDs, evt.ID)
		}

		if len(publishedIDs) > 0 {
			if err := w.repo.MarkPublished(txCtx, publishedIDs); err != nil {
				return fmt.Errorf("marking events publish: %w", err)
			}
		}

		processedCount = len(publishedIDs)
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("outbox ProcessBatch: %w", err)
	}

	return processedCount, nil
}
