package main

import (
	"net/http"
	"time"

	"nats-producer/internal/dedup"
	httpserver "nats-producer/internal/http"
	"nats-producer/internal/nats"
)

func main() {
	pub, err := nats.NewPublisher("nats://localhost:4222", "orders.created")
	if err != nil {
		panic(err)
	}

	redis := dedup.NewRedis("localhost:6379", time.Minute*10)
	local := dedup.NewLocal(time.Minute * 10)

	h := httpserver.New(pub, redis, local)

	http.HandleFunc("/send", h.HandleSend)

	println("listening on :8083")
	http.ListenAndServe(":8083", nil)
}
