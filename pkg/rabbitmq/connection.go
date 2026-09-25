package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 5
)

type Connection struct {
	url string

	attempts int
	waitTime time.Duration

	Conn *amqp.Connection
}

func NewConnection(url string, opts ...Option) *Connection {
	c := &Connection{
		url:      url,
		attempts: _defaultAttempts,
		waitTime: _defaultWaitTime,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Connection) Connect() error {
	var (
		conn *amqp.Connection
		err  error
	)
	for i := c.attempts; i > 0; i-- {
		if conn, err = amqp.Dial(c.url); err == nil {
			c.Conn = conn
			return nil
		}
		log.Printf("rabbitmq connect retry, left=%d err=%v", i-1, err)
		time.Sleep(c.waitTime)
	}
	return fmt.Errorf("rabbitmq connect exhausted retries: %v", err)
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	if c.Conn == nil {
		return nil, errors.New("rabbitmq: Channel opens before Connect")
	}

	ch, err := c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: conn.Channel: %w", err)
	}

	return ch, nil
}

func (c *Connection) Close() error {
	if c.Conn == nil {
		return nil
	}
	return c.Conn.Close()
}
