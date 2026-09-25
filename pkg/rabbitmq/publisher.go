package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Event represents a domain event to be published.
type Event struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// Publisher publishes domain events to RabbitMQ.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Close() error
}

type ExchangeBinding struct {
	Exchange     string
	ExchangeType string
}

type publisher struct {
	channel *amqp.Channel
	mu      sync.Mutex
	confirm bool
	routes  map[string]ExchangeBinding
}

func NewPublisher(conn *amqp.Connection, withConfirm bool, routes map[string]ExchangeBinding) (Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("NewPublisher - conn.Channel: %w", err)
	}

	defer func() {
		if err != nil {
			_ = ch.Close()
		}
	}()

	declared := make(map[string]bool)
	for eventType, b := range routes {
		if declared[eventType] {
			continue
		}
		if err := ch.ExchangeDeclare(b.Exchange, b.ExchangeType, true, false, false, false, nil); err != nil {
			return nil, fmt.Errorf("NewPublisher - ExchangeDeclare(%s) for event %s: %w", b.Exchange, eventType, err)
		}
		declared[eventType] = true
	}

	if withConfirm {
		if err := ch.Confirm(false); err != nil {
			return nil, fmt.Errorf("NewPublisher - ch.Confirm: %w", err)
		}
	}

	return &publisher{
		channel: ch,
		confirm: withConfirm,
		routes:  routes,
	}, nil
}

// Publish sends an event to the exchange with the event type as routing key.
func (p *publisher) Publish(ctx context.Context, event Event) error {
	binding, ok := p.routes[event.Type]
	if !ok {
		return fmt.Errorf("publisher- Publish - no exchange bound for event type %s", event.Type)
	}

	routingKey := event.Type
	if binding.ExchangeType == ExchangeTypeFanout {
		routingKey = ""
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("publisher - Publish - json.Marshal: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.confirm {
		confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
			ctx,
			binding.Exchange,
			routingKey,
			false,
			false,
			amqpPublishing(event, body),
		)
		if err != nil {
			return fmt.Errorf("publisher - Publish - PublishWithDeferredConfirmWithContext: %w", err)
		}

		ok, err := confirmation.WaitContext(ctx)
		if err != nil {
			return fmt.Errorf("publisher - Publish - confirmation.WaitContext: %w", err)
		}
		if !ok {
			return fmt.Errorf("publisher - Publish - broker nacked event %s", event.Type)
		}
		return nil
	}

	if err := p.channel.PublishWithContext(
		ctx,
		binding.Exchange,
		routingKey,
		false,
		false,
		amqpPublishing(event, body),
	); err != nil {
		return fmt.Errorf("publisher - Publish - channel.PublishWithContext: %w", err)
	}

	return nil
}

func amqpPublishing(event Event, body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    event.Timestamp,
		Type:         event.Type,
		Body:         body,
	}
}

func (p *publisher) Close() error {
	return p.channel.Close()
}
