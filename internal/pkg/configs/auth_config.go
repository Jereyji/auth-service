package configs

import "time"

type TokensExpiration struct {
	AccessTokenExpiresIn  time.Duration `yaml:"access_expiration"`
	RefreshTokenExpiresIn time.Duration `yaml:"refresh_expiration"`
	SecretKey             string        `env:"SECRET_KEY" env-required:"true"`
}

type ApplicationConfig struct {
	Tokens TokensExpiration `yaml:"tokens"`
}

type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type AuthConfig struct {
	Server      ServerConfig `yaml:"server"`
	Kafka       KafkaConfig  `yaml:"kafka"`
	Postgres    PostgresConfig
	Redis       RedisConfig
	Application ApplicationConfig
}
