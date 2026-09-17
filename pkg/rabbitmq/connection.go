package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_defaultWaitTime     = 5 * time.Second
	_defaultAttempts     = 5
	_defaultExchangeType = "topic"
)

type Connection struct {
	URL          string
	Exchange     string
	ExchangeType string

	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func newConnection(url, exchange string, opts ...Option) *Connection {
	c := &Connection{
		URL:          url,
		Exchange:     exchange,
		ExchangeType: _defaultExchangeType,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Connection) Connect() error {
	var err error
	for i := _defaultAttempts; i > 0; i-- {
		if err := c.dial(); err == nil {
			break
		}
		log.Printf("rabbitmq connect retry, left=%d err=%v", i-1, err)
		time.Sleep(_defaultWaitTime)
	}
	if err != nil {
		return fmt.Errorf("rabbitmq connect exhausted retries: %v", err)
	}
	return nil
}

func (c *Connection) dial() error {
	conn, err := amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	var ch *amqp.Channel

	defer func() {
		if err != nil {
			if ch != nil {
				_ = ch.Close()
			}
			_ = conn.Close()
		}
	}()

	ch, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		c.Exchange,
		c.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDeclare: %w", err)
	}

	c.Conn = conn
	c.Channel = ch
	return nil
}

func (c *Connection) Close() error {
	var errs []error
	if c.Channel != nil {
		if err := c.Channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if c.Conn != nil {
		if err := c.Conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
