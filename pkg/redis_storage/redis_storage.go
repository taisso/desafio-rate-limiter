package redis_storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string) *RedisStorage {
	return &RedisStorage{client: redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})}
}

func (s *RedisStorage) Incr(ctx context.Context, key string, expiration time.Duration) error {
	pipe := s.client.TxPipeline()

	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)

	return err
}

func (s *RedisStorage) Get(ctx context.Context, key string) (int, error) {
	return s.client.Get(ctx, key).Int()
}
