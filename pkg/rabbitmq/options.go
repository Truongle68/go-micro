package rabbitmq

type Option func(*Connection)

func ExchangeType(exchangeType string) Option {
	return func(c *Connection) {
		c.ExchangeType = exchangeType
	}
}
