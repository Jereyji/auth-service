package configs

import "time"

type AuthConfig struct {
	SecretKey             string        `env:"SECRET_KEY" env-required:"true"`
	AccessTokenExpiresIn  time.Duration `yaml:"access_expiration" env-required:"true"`
	RefreshTokenExpiresIn time.Duration `yaml:"refresh_expiration" env-required:"true"`
}
