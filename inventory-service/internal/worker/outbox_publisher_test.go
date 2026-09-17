package worker_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"inventory-service/internal/domain"
	"inventory-service/internal/worker"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

type mockOutboxStore struct {
	mu           sync.Mutex
	unpublished  []domain.OutboxEvent
	publishedIDs []string
	retryIDs     []string
}

func (m *mockOutboxStore) GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if limit > len(m.unpublished) {
		limit = len(m.unpublished)
	}
	return m.unpublished[:limit], nil
}

func (m *mockOutboxStore) MarkPublished(ctx context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishedIDs = append(m.publishedIDs, ids...)
	return nil
}

func (m *mockOutboxStore) IncrementRetry(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.retryIDs = append(m.retryIDs, id)
	return nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockPublisher struct {
	mu        sync.Mutex
	published []rabbitmq.Event
	failForID map[string]bool
}

func (m *mockPublisher) Publish(ctx context.Context, event rabbitmq.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var payload map[string]any
	_ = json.Unmarshal(event.Payload, &payload)

	if m.failForID != nil && m.failForID[event.Type] {
		return errors.New("simulated publish error")
	}
	m.published = append(m.published, event)
	return nil
}

func TestOutboxPublisherWorker_ProcessBatch(t *testing.T) {
	t.Run("successfully publishes and marks published", func(t *testing.T) {
		store := &mockOutboxStore{
			unpublished: []domain.OutboxEvent{
				{
					ID:            "evt-1",
					AggregateType: "purchase_order",
					AggregateID:   "po-1",
					EventType:     "GoodsReceived",
					Payload:       json.RawMessage(`{"po_code":"PO-1"}`),
					CreatedAt:     time.Now().UTC(),
					RetryCount:    0,
				},
				{
					ID:            "evt-2",
					AggregateType: "purchase_order",
					AggregateID:   "po-2",
					EventType:     "GoodsReceived",
					Payload:       json.RawMessage(`{"po_code":"PO-2"}`),
					CreatedAt:     time.Now().UTC(),
					RetryCount:    0,
				},
			},
		}

		pub := &mockPublisher{}
		l := logger.New("error")
		w := worker.NewOutboxPublisherWorker(
			store,
			&mockTransactor{},
			pub,
			l,
			worker.WithBatchSize(10),
			worker.WithMaxRetries(3),
		)

		count, err := w.ProcessBatch(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 2 {
			t.Errorf("expected processed count 2, got %d", count)
		}
		if len(store.publishedIDs) != 2 {
			t.Errorf("expected 2 published IDs, got %d", len(store.publishedIDs))
		}
		if len(pub.published) != 2 {
			t.Errorf("expected 2 published events, got %d", len(pub.published))
		}
	})

	t.Run("publishes fail triggers increment retry", func(t *testing.T) {
		store := &mockOutboxStore{
			unpublished: []domain.OutboxEvent{
				{
					ID:            "evt-fail",
					AggregateType: "purchase_order",
					AggregateID:   "po-fail",
					EventType:     "GoodsReceivedFail",
					Payload:       json.RawMessage(`{}`),
					CreatedAt:     time.Now().UTC(),
					RetryCount:    1,
				},
			},
		}

		pub := &mockPublisher{
			failForID: map[string]bool{"GoodsReceivedFail": true},
		}
		l := logger.New("error")
		w := worker.NewOutboxPublisherWorker(
			store,
			&mockTransactor{},
			pub,
			l,
		)

		count, err := w.ProcessBatch(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 published, got %d", count)
		}
		if len(store.retryIDs) != 1 || store.retryIDs[0] != "evt-fail" {
			t.Errorf("expected retry incremented for evt-fail, got %v", store.retryIDs)
		}
	})

	t.Run("skips events exceeding max retries", func(t *testing.T) {
		store := &mockOutboxStore{
			unpublished: []domain.OutboxEvent{
				{
					ID:            "evt-exceeded",
					AggregateType: "purchase_order",
					AggregateID:   "po-x",
					EventType:     "GoodsReceived",
					Payload:       json.RawMessage(`{}`),
					CreatedAt:     time.Now().UTC(),
					RetryCount:    5,
				},
			},
		}

		pub := &mockPublisher{}
		l := logger.New("error")
		w := worker.NewOutboxPublisherWorker(
			store,
			&mockTransactor{},
			pub,
			l,
			worker.WithMaxRetries(5),
		)

		count, err := w.ProcessBatch(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 processed count, got %d", count)
		}
		if len(pub.published) != 0 {
			t.Errorf("expected 0 published events, got %d", len(pub.published))
		}
	})
}
