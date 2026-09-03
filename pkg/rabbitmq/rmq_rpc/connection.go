package rmqrpc

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL      string
	WaitTime time.Duration
	Attempts int
}

type Connection struct {
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Delivery   <-chan amqp.Delivery
}

func New(consumerExchange string, cfg Config) *Connection {
	return &Connection{
		ConsumerExchange: consumerExchange,
		Config:           cfg,
	}
}

func (conn *Connection) AttemptConnect() error {
	var err error
	for i := conn.Attempts; i > 0; i-- {
		if err = conn.connect(); err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(conn.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("rmq_rpc - AttemptConnect - conn.connect: %w", err)
	}

	return nil
}

func (conn *Connection) connect() error {
	var err error

	conn.Connection, err = amqp.Dial(conn.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	conn.Channel, err = conn.Connection.Channel()
	if err != nil {
		return fmt.Errorf("conn.Connection.Channel: %w", err)
	}

	err = conn.Channel.ExchangeDeclare(conn.ConsumerExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("conn.Channel.ExchangeDeclare: %w", err)
	}

	queue, err := conn.Channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return fmt.Errorf("conn.Channel.QueueDeclare: %w", err)
	}

	err = conn.Channel.QueueBind(queue.Name, "", conn.ConsumerExchange, false, nil)
	if err != nil {
		return fmt.Errorf("conn.Channel.QueueBind: %w", err)
	}

	conn.Delivery, err = conn.Channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("conn.Channel.Consume: %w", err)
	}

	return nil
}
