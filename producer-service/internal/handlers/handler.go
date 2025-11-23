package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"producer-service/internal/dedup"
	"producer-service/internal/kafka"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	dedup *dedup.Service
	prod  *kafka.Producer
}

type Message struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func New(d *dedup.Service, p *kafka.Producer) *Handler {
	return &Handler{dedup: d, prod: p}
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	seen, _ := h.dedup.Seen(ctx, msg.ID)
	if seen {
		http.Error(w, "duplicate", 409)
		log.Log().Msgf("Duplicated: %s", msg.ID)
		return
	}

	h.dedup.Mark(ctx, msg.ID)

	if err := h.prod.Send(ctx, []byte(msg.ID), []byte(msg.Payload)); err != nil {
		log.Error().Err(err).Msg("kafka send")
		http.Error(w, "kafka error", 500)
		return
	}
	log.Log().Msgf("kafka sent: %s", msg.ID)
	w.Write([]byte("Event processed to kafka"))
}
