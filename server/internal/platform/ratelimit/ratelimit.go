package ratelimit

import (
	"sync"
	"time"
)

type Memory struct {
	mu      sync.Mutex
	buckets map[string]*window
	now     func() time.Time
	stop    chan struct{}
}

type window struct {
	count   int
	resetAt time.Time
}

func NewMemory() *Memory {
	m := &Memory{
		buckets: make(map[string]*window),
		now:     time.Now,
		stop:    make(chan struct{}),
	}

	go m.janitor()

	return m
}

func (m *Memory) Allow(key string, limit int, win time.Duration) (allowed bool, retryAfter time.Duration) {
	now := m.now()

	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.buckets[key]
	if !ok || now.After(w.resetAt) {
		m.buckets[key] = &window{
			count:   1,
			resetAt: now.Add(win),
		}
		return true, 0
	}

	if w.count >= limit {
		return false, w.resetAt.Sub(now)
	}

	w.count++
	return true, 0
}

func (m *Memory) janitor() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		select {
		case <-m.stop:
			return
		case <-t.C:
			now := m.now()
			m.mu.Lock()
			for k, w := range m.buckets {
				if now.After(w.resetAt) {
					delete(m.buckets, k)
				}
			}
			m.mu.Unlock()
		}
	}
}

func (m *Memory) Close() { close(m.stop) }
