package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer interface {
	Start(ctx context.Context) error
	Close() error
}

type consumer struct {
	conn       *Connection
	queueName  string
	routingKey string
	handler    HandlerFunc
	deliveries <-chan amqp.Delivery
}

func NewConsumer(
	url string,
	exchange string,
	queueName string,
	routingKey string,
	handler HandlerFunc,
	opts ...Option,
) (_ Consumer, err error) {
	c := newConnection(url, exchange, opts...)
	if err = c.Connect(); err != nil {
		return nil, fmt.Errorf("NewConsumer - c.Connect: %w", err)
	}

	defer func() {
		if err != nil {
			_ = c.Conn.Close()
		}
	}()

	q, err := c.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewConsumer - QueueDeclare: %w", err)
	}

	err = c.Channel.QueueBind(
		q.Name,
		routingKey,
		c.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewConsumer - QueueBind: %w", err)
	}

	deliveries, err := c.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewConsumer - Consume: %w", err)
	}

	return &consumer{
		conn:       c,
		queueName:  q.Name,
		routingKey: routingKey,
		handler:    handler,
		deliveries: deliveries,
	}, nil
}

func (c *consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, open := <-c.deliveries:
			if !open {
				return fmt.Errorf("consumer: delivery channel closed")
			}
			if err := c.handleOne(ctx, d); err != nil {
				_ = d.Nack(false, true)
				continue
			}
			_ = d.Ack(false)
		}
	}
}

func (c *consumer) handleOne(ctx context.Context, deliveries amqp.Delivery) error {
	var event Event
	if err := json.Unmarshal(deliveries.Body, &event); err != nil {
		return fmt.Errorf("unmarshaling event: %w", err)
	}
	return c.handler(ctx, event)
}

func (c *consumer) Close() error {
	return c.conn.Close()
}
