package configs

import "fmt"

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     int    `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB" env-required:"true"`
}

func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
