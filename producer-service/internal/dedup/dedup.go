package dedup

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	client *redis.Client
}

func New(addr string) *Service {
	return &Service{client: redis.NewClient(&redis.Options{Addr: addr})}
}

func (s *Service) Seen(ctx context.Context, id string) (bool, error) {
	exists, err := s.client.Exists(ctx, "msg:"+id).Result()
	return exists == 1, err
}

func (s *Service) Mark(ctx context.Context, id string) error {
	return s.client.Set(ctx, "msg:"+id, 1, 24*time.Hour).Err()
}
