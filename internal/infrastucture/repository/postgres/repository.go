package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type EstateRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func NewEstateRepository(db *pgxpool.Pool) *EstateRepository {
	return &EstateRepository{
		db: db,
	}
}
