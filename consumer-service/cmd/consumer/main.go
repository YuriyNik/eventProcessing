package main

import (
	//"consumer-service/internal/config"
	"consumer-service/internal/kafka"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	consumer := kafka.NewConsumer(
		[]string{"localhost:9092"},
		"events",
		"batch-consumer-group",
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("Starting batch consumer…")
	consumer.Run(ctx)
}
