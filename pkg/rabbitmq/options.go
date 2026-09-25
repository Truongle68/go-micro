package rabbitmq

import "time"

type Option func(*Connection)

func Attempts(n int) Option {
	return func(c *Connection) {
		if n > 0 {
			c.attempts = n
		}
	}
}

func WaitTime(d time.Duration) Option {
	return func(c *Connection) {
		if d > 0 {
			c.waitTime = d
		}
	}
}
