package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/redis"
	tr_manager "github.com/Jereyji/auth-service/internal/pkg/transaction_manager"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	cashingTime      = 15 * time.Minute
	userCacheKeyText = string("user:%s")
	rtCacheKeyText   = string("refresh_token:%s")
)

type AuthRepository struct {
	trm         tr_manager.TransactionManager
	redisClient *redis.RedisClient
}

func NewAuthRepository(trm tr_manager.TransactionManager, redisClient *redis.RedisClient) *AuthRepository {
	return &AuthRepository{
		trm:         trm,
		redisClient: redisClient,
	}
}

func ifUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}

	return false
}

func formatCacheKey(mainText string, value any) string {
	return fmt.Sprintf(mainText, value)
}
