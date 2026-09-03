package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_defaultExchange = "inventory.events"
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
)

// Event represents a domain event to be published.
type Event struct {
	Type      string    `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

// Publisher publishes domain events to RabbitMQ.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Close() error
}

type rabbitPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	url      string
}

// New creates a new RabbitMQ event publisher.
// It connects to RabbitMQ and declares a topic exchange.
func New(url string, opts ...Option) (Publisher, error) {
	p := &rabbitPublisher{
		exchange: _defaultExchange,
		url:      url,
	}

	for _, opt := range opts {
		opt(p)
	}

	if err := p.connect(); err != nil {
		return nil, fmt.Errorf("publisher.New - connect: %w", err)
	}

	return p, nil
}

// Option configures the Publisher.
type Option func(*rabbitPublisher)

// Exchange sets the exchange name.
func Exchange(name string) Option {
	return func(p *rabbitPublisher) {
		p.exchange = name
	}
}

func (p *rabbitPublisher) connect() error {
	var err error
	for i := _defaultAttempts; i > 0; i-- {
		if err = p.dial(); err == nil {
			return nil
		}
		log.Printf("publisher - connect - attempting, retries left: %d", i-1)
		time.Sleep(_defaultWaitTime)
	}
	return fmt.Errorf("publisher - connect - exhausted retries: %w", err)
}

func (p *rabbitPublisher) dial() error {
	var err error

	p.conn, err = amqp.Dial(p.url)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}

	// Declare a topic exchange for routing events by type
	err = p.channel.ExchangeDeclare(
		p.exchange,
		"topic", // topic exchange allows routing key patterns
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("channel.ExchangeDeclare: %w", err)
	}

	return nil
}

// Publish sends an event to the exchange with the event type as routing key.
func (p *rabbitPublisher) Publish(ctx context.Context, event Event) error {
	if p.channel == nil {
		return errors.New("publisher - Publish - channel is nil")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("publisher - Publish - json.Marshal: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,
		event.Type, // routing key = event type (e.g. "stock.reserved")
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.Timestamp,
			Type:         event.Type,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publisher - Publish - channel.PublishWithContext: %w", err)
	}

	return nil
}

// Close tears down the channel and connection.
func (p *rabbitPublisher) Close() error {
	var errs []error

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
