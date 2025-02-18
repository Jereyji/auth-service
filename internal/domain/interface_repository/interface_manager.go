package repos

import "context"

type TransactionManagerI interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}