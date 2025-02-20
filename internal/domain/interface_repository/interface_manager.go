package repos

import "context"

type TransactionManagerI interface {
	WithTransaction(ctx context.Context, f func(ctx context.Context) error) error
}