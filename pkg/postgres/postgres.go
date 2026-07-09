package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize = 3
	_defaultConnAttempt = 5
	_defaultConnTimeout = 5 * time.Second
)

type Postgres struct {
	maxPoolSize int
	connAttempt int
	connTimeout time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize: _defaultMaxPoolSize,
		connAttempt: _defaultConnAttempt,
		connTimeout: _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	pgConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	pgConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempt > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
			err = pg.Pool.Ping(ctx)
			cancel()
		}

		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d\n", pg.connAttempt)
		time.Sleep(pg.connTimeout)
		pg.connAttempt--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return pg, nil
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
