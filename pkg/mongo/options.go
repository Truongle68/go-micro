package mongo

import "time"

type Option func(*Mongo)

func MaxPoolSize(size int) Option {
	return func(m *Mongo) {
		if size > 0 {
			m.maxPoolSize = size
		}
	}
}

func ConnAttempt(attempt int) Option {
	return func(m *Mongo) {
		if attempt > 0 {
			m.connAttempt = attempt
		}
	}
}

func ConnTimeout(t time.Duration) Option {
	return func(m *Mongo) {
		if t > 0 {
			m.connTimeout = t
		}
	}
}
