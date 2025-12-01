package dedup

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDedup struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(addr string, ttl time.Duration) *RedisDedup {
	return &RedisDedup{
		client: redis.NewClient(&redis.Options{Addr: addr}),
		ttl:    ttl,
	}
}

func (d *RedisDedup) Seen(ctx context.Context, id string) (bool, error) {
	ok, err := d.client.SetNX(ctx, id, 1, d.ttl).Result()
	if err != nil {
		return false, err
	}
	return !ok, nil
}
