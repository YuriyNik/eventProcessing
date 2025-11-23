package main

import (
	"producer-service/internal/config"
	"producer-service/internal/dedup"
	"producer-service/internal/handlers"
	"producer-service/internal/httpserver"
	"producer-service/internal/kafka"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()

	dedupSvc := dedup.New(cfg.RedisAddr)
	producer := kafka.New(cfg.KafkaBrokers, "events")
	h := handlers.New(dedupSvc, producer)

	r := mux.NewRouter()
	r.HandleFunc("/send", h.Send).Methods("POST")

	s := httpserver.New(cfg.Port, r)

	log.Info().Msg("Producer service running")
	s.Start()
}
