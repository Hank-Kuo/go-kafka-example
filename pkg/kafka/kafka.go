package kafka

import (
	"errors"

	"github.com/Hank-Kuo/go-kafka-example/config"

	"github.com/segmentio/kafka-go"
)

func NewWriter(cfg config.KafkaConfig, topicName string) (*kafka.Writer, error) {
	topic, ok := cfg.Topics[topicName]

	if !ok {
		return nil, errors.New("not found topic")
	}

	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    topic.Name,
		Balancer: &kafka.LeastBytes{},
	}, nil
}

func NewReader(cfg config.KafkaConfig, topicName string) (*kafka.Reader, error) {
	topic, ok := cfg.Topics[topicName]
	if !ok {
		return nil, errors.New("not found topic")
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    topic.Name,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		GroupID:  cfg.GroupID,
	}), nil
}
