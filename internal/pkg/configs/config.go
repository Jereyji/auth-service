package configs

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"postgres"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"auth_db"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type AuthConfig struct {
	SecretKey             string        `env:"SECRET_KEY" env-required:"true"`
	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRATION" env-default:"30m"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRATION" env-default:"48h"`
}

type Config struct {
	Database    DatabaseConfig // Сделали экспортируемым
	DatabaseURL string         `env:"-"`
	AuthService AuthConfig     // Сделали экспортируемым
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	// Формируем URL для базы данных
	cfg.DatabaseURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	return &cfg, nil
}
