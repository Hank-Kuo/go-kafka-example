package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewWriter(url string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewReader(url []string, topic string) *kafka.Reader {

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  url,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		GroupID:  "customer-1",
	})
}
