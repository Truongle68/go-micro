package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"inventory-service/internal/domain"
	invpg "inventory-service/pkg/postgres"

	"github.com/google/uuid"
)

type OutboxRepo struct {
	db *sql.DB
}

type OutboxRepository = OutboxRepo

func NewOutboxRepo(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{db: db}
}

func NewOutboxRepository(db *sql.DB) *OutboxRepo {
	return NewOutboxRepo(db)
}

func (r *OutboxRepo) Create(ctx context.Context, event *domain.OutboxEvent) error {
	executor := invpg.GetExecutor(ctx, r.db)

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO outbox (
			id, aggregate_type, aggregate_id, event_type, payload, created_at, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := executor.ExecContext(ctx, query,
		event.ID,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.Payload,
		event.CreatedAt,
		event.RetryCount,
	)
	if err != nil {
		return fmt.Errorf("OutboxRepo.Create: %w", err)
	}
	return nil
}

func (r *OutboxRepo) GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	executor := invpg.GetExecutor(ctx, r.db)
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, created_at, published_at, retry_count
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	rows, err := executor.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("OutboxRepo.GetUnpublished - query: %w", err)
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var evt domain.OutboxEvent
		var payload []byte
		if err := rows.Scan(
			&evt.ID,
			&evt.AggregateType,
			&evt.AggregateID,
			&evt.EventType,
			&payload,
			&evt.CreatedAt,
			&evt.PublishedAt,
			&evt.RetryCount,
		); err != nil {
			return nil, fmt.Errorf("OutboxRepo.GetUnpublished - scan: %w", err)
		}
		evt.Payload = json.RawMessage(payload)
		events = append(events, evt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("OutboxRepo.GetUnpublished - rows: %w", err)
	}
	return events, nil
}

func (r *OutboxRepo) MarkPublished(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	executor := invpg.GetExecutor(ctx, r.db)

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
		return fmt.Errorf("OutboxRepo.MarkPublished: %w", err)
	}
	return nil
}

func (r *OutboxRepo) IncrementRetry(ctx context.Context, id string) error {
	executor := invpg.GetExecutor(ctx, r.db)

	query := `
		UPDATE outbox
		SET retry_count = retry_count + 1
		WHERE id = $1`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("OutboxRepo.IncrementRetry: %w", err)
	}
	return nil
}
