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
	dedup         *dedup.Service
	localFallback *dedup.LocalDedup
	prod          *kafka.Producer
}

type Message struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func New(d *dedup.Service, l *dedup.LocalDedup, p *kafka.Producer) *Handler {
	return &Handler{dedup: d, localFallback: l, prod: p}
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	seen, err := h.dedup.Seen(ctx, msg.ID)
	if err != nil {
		// Redis not available -  fallback
		//log.Warn().Err(err).Msg("redis unavailable, using local fallback")

		if h.localFallback.Seen(msg.ID) {
			http.Error(w, "duplicate (local)", http.StatusConflict)
			log.Warn().Err(err).Msg("redis unavailable, using local fallback - duplicate found for id " + msg.ID)
			return
		}
	} else {
		// Redis available - regular flow
		if seen {
			http.Error(w, "duplicate", http.StatusConflict)
			log.Warn().Err(err).Msg("redis available, duplicate found for id " + msg.ID)
			return
		}
		// store id at redis
		if err := h.dedup.Mark(ctx, msg.ID); err != nil {
			// if some error occurs - fallback to local
			log.Warn().Err(err).Msg("failed to mark in redis, using local fallback")
			if h.localFallback.Seen(msg.ID) {
				http.Error(w, "duplicate (local)", http.StatusConflict)
				return
			}
		}
	}

	if err := h.prod.Send(ctx, []byte(msg.ID), []byte(msg.Payload)); err != nil {
		log.Error().Err(err).Msg("kafka send")
		http.Error(w, "kafka error", 500)
		return
	}
	log.Log().Msgf("kafka sent: %s", msg.ID)
	w.Write([]byte("Event processed to kafka"))
}
