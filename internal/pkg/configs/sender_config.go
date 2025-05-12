package configs

type SenderConfig struct {
	EnvMode string `yaml:"env"`

	Kafka KafkaConfig `yaml:"kafka"`

	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `env:"EMAIL" env-required:"true"`
	Password string `env:"EMAIL_PASSWORD" env-required:"true"`
}
