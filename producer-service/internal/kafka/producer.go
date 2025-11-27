package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func New(brokers []string, topic string) *Producer {
	addr := "127.0.0.1:9092"
	if len(brokers) > 0 {
		addr = brokers[0]
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchSize:    100,
		BatchTimeout: 200 * time.Millisecond,
	}

	return &Producer{writer: w}
}

func (p *Producer) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

// Close закрывает writer, освобождая ресурсы.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Shared — пакетный синглтон для повторного использования writer по всему приложению.
var Shared *Producer

// InitShared инициализирует Shared один раз.
func InitShared(brokers []string, topic string) {
	if Shared != nil {
		return
	}
	Shared = New(brokers, topic)
}
