package configs

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"consumer_group"`
}
