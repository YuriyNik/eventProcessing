package main

import (
	"producer-service/internal/config"
	"producer-service/internal/dedup"
	"producer-service/internal/handlers"
	"producer-service/internal/httpserver"
	"producer-service/internal/kafka"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	localDedup := dedup.NewLocalDedup(30 * time.Second) // TTL 30s

	dedupSvc := dedup.New(cfg.RedisAddr)
	producer := kafka.New(cfg.KafkaBrokers, "events")
	h := handlers.New(dedupSvc, localDedup, producer)

	r := mux.NewRouter()
	r.HandleFunc("/send", h.Send).Methods("POST")

	s := httpserver.New(cfg.Port, r)

	log.Info().Msg("Producer service running")
	s.Start()
}
