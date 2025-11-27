package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"producer-service/internal/config"
	"producer-service/internal/dedup"
	"producer-service/internal/handlers"
	"producer-service/internal/kafka"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	cfg := config.Load()

	brokers := []string{getEnv("KAFKA_BROKERS", "127.0.0.1:9092")}
	topic := getEnv("KAFKA_TOPIC", "events")

	// Инициализация общего продюсера
	kafka.InitShared(brokers, topic)

	localDedup := dedup.NewLocalDedup(30 * time.Second) // TTL 30s

	dedupSvc := dedup.New(cfg.RedisAddr)
	//producer := kafka.New(cfg.KafkaBrokers, "events")
	h := handlers.New(dedupSvc, localDedup)

	http.HandleFunc("/send", h.Send)
	srv := &http.Server{
		Addr:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
