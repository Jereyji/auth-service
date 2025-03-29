package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *configs.RedisConfig) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Address(),
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

func (r *RedisClient) Close(ctx context.Context) error {
	return r.client.Close()
}

func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = r.client.Set(ctx, key, valueJSON, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string, value any) error {
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(cachedValue), value); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}
