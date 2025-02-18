package transaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	tx *pgxpool.Tx
}

func NewTransaction(tx *pgxpool.Tx) *Transaction {
	return &Transaction{
		tx: tx,
	}
}

func (t *Transaction) Transaction() interface{} {
	return t.tx
}

func (t *Transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}