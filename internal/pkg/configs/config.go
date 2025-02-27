package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "deployments/.env"

type Config struct {
	Database    DatabaseConfig
	AuthService AuthConfig   `yaml:"tokens"`
	Server      ServerConfig `yaml:"server"`
}

func NewConfig(configPath string) (*Config, error) {
	var cfg Config

	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
