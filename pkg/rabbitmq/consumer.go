package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer struct {
	channel    *amqp.Channel
	queueName  string
	handler    HandlerFunc
	deliveries <-chan amqp.Delivery
}

func NewConsumer(
	conn *amqp.Connection,
	exchange, queueName, routingKey string,
	handler HandlerFunc,
) (_ *Consumer, err error) {
	return newConsumer(conn, exchange, ExchangeTypeTopic, queueName, routingKey, handler)
}

func NewFanoutConsumer(conn *amqp.Connection, exchange, queueName string, handler HandlerFunc) (*Consumer, error) {
	return newConsumer(conn, exchange, ExchangeTypeFanout, queueName, "", handler)
}

func newConsumer(
	conn *amqp.Connection,
	exchange, exchangeType, queueName, routingKey string,
	handler HandlerFunc,
) (_ *Consumer, err error) {
	ch, err := conn.Channel()
	defer func() {
		if err != nil {
			_ = ch.Close()
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("conn.Channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("NewConsumer - ExchangeDeclare: %w", err)
	}

	q, err := ch.QueueDeclare(
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

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("NewConsumer - QueueBind: %w", err)
	}

	deliveries, err := ch.Consume(
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

	return &Consumer{
		channel:    ch,
		queueName:  q.Name,
		handler:    handler,
		deliveries: deliveries,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, open := <-c.deliveries:
			if !open {
				return fmt.Errorf("Consumer: delivery channel closed")
			}
			if err := c.handleOne(ctx, d); err != nil {
				_ = d.Nack(false, true)
				continue
			}
			_ = d.Ack(false)
		}
	}
}

func (c *Consumer) handleOne(ctx context.Context, deliveries amqp.Delivery) error {
	var event Event
	if err := json.Unmarshal(deliveries.Body, &event); err != nil {
		return fmt.Errorf("unmarshaling event: %w", err)
	}
	return c.handler(ctx, event)
}

func (c *Consumer) Close() error {
	return c.channel.Close()
}
