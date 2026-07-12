package postgres

import "time"

type Option func(*Postgres)

func MaxOpenConns(size int) Option {
	return func(p *Postgres) { p.maxOpenConns = size }
}

func MaxIdleConns(size int) Option {
	return func(p *Postgres) { p.maxIdleConns = size }
}

func ConnMaxLifetime(d time.Duration) Option {
	return func(p *Postgres) { p.connMaxLifetime = d }
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgres) { p.connAttempt = attempts }
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) { p.connTimeout = timeout }
}
