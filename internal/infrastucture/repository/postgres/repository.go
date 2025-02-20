package repository

import (
	"errors"

	trm "github.com/Jereyji/auth-service.git/pkg/transaction_manager"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository struct {
	txm trm.TransactionManager
}

func NewAuthRepository(txm trm.TransactionManager) *AuthRepository {
	return &AuthRepository{
		txm: txm,
	}
}

func ifUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}

	return false
}
