package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	_defaultMaxOpenConns    = 3
	_defaultMaxIdleConns    = 3
	_defaultConnMaxLifetime = 30 * time.Minute
	_defaultConnAttempt     = 5
	_defaultConnTimeout     = 5 * time.Second
)

type Postgres struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connAttempt     int
	connTimeout     time.Duration

	DB *sql.DB
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxOpenConns:    _defaultMaxOpenConns,
		maxIdleConns:    _defaultMaxIdleConns,
		connMaxLifetime: _defaultConnMaxLifetime,
		connAttempt:     _defaultConnAttempt,
		connTimeout:     _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - sql.Open: %w", err)
	}

	db.SetMaxOpenConns(pg.maxOpenConns)
	db.SetMaxIdleConns(pg.maxIdleConns)
	db.SetConnMaxLifetime(pg.connMaxLifetime)

	pg.DB = db

	for pg.connAttempt > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
		err = pg.DB.PingContext(ctx)
		cancel()

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
	if pg.DB != nil {
		pg.DB.Close()
	}
}
