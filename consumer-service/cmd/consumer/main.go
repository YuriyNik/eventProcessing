package main

import (
	"consumer-service/internal/config"
	"consumer-service/internal/kafka"
	"context"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	consumer := kafka.New(cfg.KafkaBrokers, cfg.Topic, cfg.GroupID)

	log.Info().Msg("consumer started")

	consumer.Run(context.Background(), func(key, value []byte) {
		log.Info().Msgf("Processing message... received key=%s value=%s", key, value)
	})
}
