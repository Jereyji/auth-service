package repos

import "context"

type ITransactionManager interface {
	WithTransaction(ctx context.Context, f func(ctx context.Context) error) error
}
