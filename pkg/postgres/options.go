package postgres

import "time"

type Option func(*Postgres)

func MaxSizePool(size int) Option {
	return func(p *Postgres) {
		if size > 0 {
			p.maxPoolSize = size
		}
	}
}

func ConnAttempt(attempt int) Option {
	return func(p *Postgres) {
		if attempt > 0 {
			p.connAttempt = attempt
		}
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		if timeout > 0 {
			p.connTimeout = timeout
		}
	}
}
