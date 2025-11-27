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
}

type Message struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func New(d *dedup.Service, l *dedup.LocalDedup) *Handler {
	return &Handler{dedup: d, localFallback: l}
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
		if h.localFallback.Seen(msg.ID) {
			http.Error(w, "duplicate (local)", http.StatusConflict)
			log.Warn().Err(err).Msg("redis unavailable, using local fallback - duplicate found for id " + msg.ID)
			return
		}
	} else {
		if seen {
			http.Error(w, "duplicate", http.StatusConflict)
			log.Warn().Err(err).Msg("redis available, duplicate found for id " + msg.ID)
			return
		}
		if err := h.dedup.Mark(ctx, msg.ID); err != nil {
			log.Warn().Err(err).Msg("failed to mark in redis, using local fallback")
			if h.localFallback.Seen(msg.ID) {
				http.Error(w, "duplicate (local)", http.StatusConflict)
				return
			}
		}
	}

	if kafka.Shared == nil {
		log.Error().Msg("kafka shared producer is not initialized")
		http.Error(w, "kafka not initialized", http.StatusInternalServerError)
		return
	}

	if err := kafka.Shared.Send(ctx, []byte(msg.ID), []byte(msg.Payload)); err != nil {
		log.Error().Err(err).Msg("kafka send")
		http.Error(w, "kafka error", 500)
		return
	}

	w.Write([]byte("Event processed to kafka"))
}
