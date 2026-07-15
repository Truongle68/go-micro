package postgres

import (
	"context"
	"database/sql"

	"github.com/GoProOrg/core-go-pkg/database"
)

type sqlTxKey struct{}

type PostgresTransactor struct {
	db *sql.DB
}

func NewPostgresTransactor(db *sql.DB) *PostgresTransactor {
	return &PostgresTransactor{
		db: db,
	}
}

func (t *PostgresTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return database.GenericWithTransaction(ctx, t.db, func(ctx context.Context, tx *sql.Tx) context.Context {
		return context.WithValue(ctx, sqlTxKey{}, tx)
	}, fn)
}

func GetExecutor(ctx context.Context, fallback *sql.DB) database.DBExecutor {
	if tx, ok := ctx.Value(sqlTxKey{}).(*sql.Tx); ok {
		return tx
	}
	return fallback
}
