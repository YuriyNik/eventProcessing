package kafka

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func New(brokers, topic, group string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokers},
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (c *Consumer) Run(ctx context.Context, handler func(key, val []byte)) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("consumer read error")
			continue
		}
		handler(m.Key, m.Value)
	}
}
