package configs

import "time"

type GinConfig struct {
	Mode      string   `yaml:"gin_mode"`
	SkipPaths []string `yaml:"skip_paths"`
}

type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type TokensConfig struct {
	AccessTokenExpiresIn  time.Duration `yaml:"access_expiration"`
	RefreshTokenExpiresIn time.Duration `yaml:"refresh_expiration"`
	SecretKey             string        `env:"SECRET_KEY" env-required:"true"`
}

type AuthConfig struct {
	Gin      GinConfig    `yaml:"gin"`
	Kafka    KafkaConfig  `yaml:"kafka"`
	Server   ServerConfig `yaml:"server"`
	Tokens   TokensConfig `yaml:"tokens"`
	Postgres PostgresConfig
	Redis    RedisConfig
}
