package dedup

import (
	"sync"
	"time"
)

type LocalDedup struct {
	mu   sync.Mutex
	seen map[string]time.Time
	ttl  time.Duration
}

func NewLocal(ttl time.Duration) *LocalDedup {
	return &LocalDedup{
		seen: make(map[string]time.Time),
		ttl:  ttl,
	}
}

func (d *LocalDedup) Seen(id string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	if exp, ok := d.seen[id]; ok && exp.After(now) {
		return true
	}

	d.seen[id] = now.Add(d.ttl)
	return false
}
