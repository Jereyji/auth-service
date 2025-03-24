package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func NewConfig[T any](cfg *T, configPath, envPath string) error {
	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return err
	}

	return nil
}
