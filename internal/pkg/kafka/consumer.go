package kafka

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	handler  sarama.ConsumerGroupHandler
}

func NewKafkaConsumer(brokers []string, groupID, topic string, handler sarama.ConsumerGroupHandler) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}, nil
}

func (kc *KafkaConsumer) StartListening(ctx context.Context) {
	go func() {
		for {
			if err := kc.consumer.Consume(ctx, []string{kc.topic}, kc.handler); err != nil {
				log.Printf("Ошибка при чтении Kafka-сообщений: %v", err)
			}
		}
	}()
	log.Println("Kafka Consumer запущен и слушает топик:", kc.topic)
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
