package postgres

import "context"

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
