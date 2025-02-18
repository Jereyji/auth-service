package repos

import "context"

type TransactionI interface {
	Transaction() interface{}
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
