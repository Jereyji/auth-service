package configs

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"postgres"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"auth_db"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}
