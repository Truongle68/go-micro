package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"order-service/internal/domain"
	orderpg "order-service/pkg/postgres"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OutboxRepo struct {
	db *sql.DB
}

func NewOutboxRepo(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{
		db: db,
	}
}

func (r *OutboxRepo) Create(ctx context.Context, event *domain.OutboxEvent) error {
	executor := orderpg.GetExecutor(ctx, r.db)
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO outbox (
			id, aggregate_id, aggregate_type, event_type, payload, created_at, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := executor.ExecContext(ctx, query, event.ID, event.AggregateID, event.AggregateType, event.EventType, event.Payload, event.CreatedAt, event.RetryCount)
	if err != nil {
		return fmt.Errorf("[Order Service] inserting outbox: %w", err)
	}

	return nil
}
func (r *OutboxRepo) GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	executor := orderpg.GetExecutor(ctx, r.db)
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, aggregate_id, aggregate_type, event_type, payload, created_at, published_at, retry_count
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := executor.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("[Order Service] querying outbox: %w", err)
	}
	defer rows.Close()

	var outboxEvents []domain.OutboxEvent
	for rows.Next() {
		var event domain.OutboxEvent
		var payload []byte
		if err := rows.Scan(
			&event.ID,
			&event.AggregateID,
			&event.AggregateType,
			&event.EventType,
			&payload,
			&event.CreatedAt,
			&event.PublishedAt,
			&event.RetryCount,
		); err != nil {
			return nil, fmt.Errorf("[Order Service] scanning: %w", err)
		}

		event.Payload = json.RawMessage(payload)
		outboxEvents = append(outboxEvents, event)
	}

	return outboxEvents, nil
}
func (r *OutboxRepo) MarkPublished(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	executor := orderpg.GetExecutor(ctx, r.db)

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE outbox
		SET published_at = NOW()
		WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	if _, err := executor.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("[Order Service] marking published: %w", err)
	}

	return nil
}
func (r *OutboxRepo) IncrementRetry(ctx context.Context, id string) error {
	executor := orderpg.GetExecutor(ctx, r.db)

	query := `
		UPDATE outbox
		SET retry_count = retry_count + 1
		WHERE id = $1`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("[Order Service] increasing retry: %w", err)
	}

	return nil
}
