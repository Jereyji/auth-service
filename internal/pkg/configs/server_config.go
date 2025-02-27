package configs

import "time"

type ServerConfig struct {
	Address      string        `yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-required:"true"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-required:"true"`
}
