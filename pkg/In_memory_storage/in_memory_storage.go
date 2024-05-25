package in_memory_storage

import (
	"context"
	"sync"
	"time"
)

type InMemoryStorage struct {
	mu     sync.Mutex
	values map[string]rateLimitCounter
}

type rateLimitCounter struct {
	Count  int
	Expiry time.Time
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		values: make(map[string]rateLimitCounter),
	}
}

func (s *InMemoryStorage) Incr(ctx context.Context, key string, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter, ok := s.values[key]
	if !ok || time.Now().After(counter.Expiry) {
		counter = rateLimitCounter{
			Count:  1,
			Expiry: time.Now().Add(expiration),
		}
		s.values[key] = counter
		return nil
	}

	counter.Count++
	s.values[key] = counter
	return nil
}

func (s *InMemoryStorage) Get(ctx context.Context, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter, ok := s.values[key]
	if !ok || time.Now().After(counter.Expiry) {
		return 0, nil
	}

	return counter.Count, nil
}
