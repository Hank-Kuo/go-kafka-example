package kafka

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	Url      string
	MinBytes int
	MaxBytes int
}

func (k *Kafka) NewWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.Url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *Kafka) NewReader(topic string) *kafka.Reader {
	brokers := strings.Split(k.Url, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: k.MinBytes, // 10e3, // 10KB
		MaxBytes: k.MaxBytes, //10e6, // 10MB
	})
}
