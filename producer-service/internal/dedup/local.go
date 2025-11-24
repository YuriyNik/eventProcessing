package dedup

import (
	"sync"
	"time"
)

type LocalDedup struct {
	mu   sync.Mutex
	data map[string]time.Time
	ttl  time.Duration
}

func NewLocalDedup(ttl time.Duration) *LocalDedup {
	l := &LocalDedup{
		data: make(map[string]time.Time),
		ttl:  ttl,
	}

	// start background cleanup goroutine
	go l.cleanup()

	return l
}

func (l *LocalDedup) Seen(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	t, ok := l.data[id]
	if !ok {
		l.data[id] = time.Now()
		return false
	}

	if time.Since(t) > l.ttl {
		l.data[id] = time.Now()
		return false
	}

	return true
}

func (l *LocalDedup) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for k, t := range l.data {
			if now.Sub(t) > l.ttl {
				delete(l.data, k)
			}
		}
		l.mu.Unlock()
	}
}
