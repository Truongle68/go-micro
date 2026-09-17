package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type publisher struct {
	conn *Connection
}

// New creates a new RabbitMQ event publisher.
// It connects to RabbitMQ and declares a topic exchange.
func NewPublisher(url, exchange string, opts ...Option) (Publisher, error) {
	c := newConnection(url, exchange, opts...)

	if err := c.Connect(); err != nil {
		return nil, fmt.Errorf("NewPublisher - c.Connect: %w", err)
	}

	return &publisher{
		conn: c,
	}, nil
}

// Publish sends an event to the exchange with the event type as routing key.
func (p *publisher) Publish(ctx context.Context, event Event) error {
	if p.conn.Channel == nil {
		return errors.New("publisher - Publish - channel is nil")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("publisher - Publish - json.Marshal: %w", err)
	}

	err = p.conn.Channel.PublishWithContext(
		ctx,
		p.conn.Exchange,
		event.Type,
		false,
		false,
		amqpPublishing(event, body),
	)
	if err != nil {
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
	return p.conn.Close()
}
