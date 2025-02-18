package transaction

import "github.com/jackc/pgx/v5"

type Transaction struct {
	Tx pgx.Tx
}