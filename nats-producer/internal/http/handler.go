package httpserver

import (
	"encoding/json"
	"net/http"

	"nats-producer/internal/dedup"
	"nats-producer/internal/nats"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	pub   *nats.Publisher
	redis *dedup.RedisDedup
	local *dedup.LocalDedup
}

type Message struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func New(pub *nats.Publisher, redis *dedup.RedisDedup, local *dedup.LocalDedup) *Handler {
	return &Handler{pub: pub, redis: redis, local: local}
}

func (h *Handler) HandleSend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	// try Redis
	seen, err := h.redis.Seen(ctx, msg.ID)
	if err != nil {
		log.Warn().Err(err).Msg("redis unavailable – using local fallback")

		// fallback
		if h.local.Seen(msg.ID) {
			http.Error(w, "duplicate (local)", 409)
			return
		}
	} else if seen {
		http.Error(w, "duplicate", 409)
		return
	}

	// publish to NATS
	if err := h.pub.Publish(ctx, msg); err != nil {
		log.Error().Err(err).Msg("publish failed")
		http.Error(w, "publish failed", 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
