package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ctx    context.Context
}

func New(addr string) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Store{
		client: client,
		ctx:    ctx,
	}, nil
}

func (s *Store) Close() error {
	return s.client.Close()
}

func (s *Store) rateLimitKey(userID string) string {
	return fmt.Sprintf("ratelimit:%s", userID)
}

func (s *Store) IncrementAndGet(userID string, windowSize int64) (int64, error) {
	key := s.rateLimitKey(userID)

	count, err := s.client.Incr(s.ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		s.client.Expire(s.ctx, key, time.Duration(windowSize)*time.Second)
	}

	return count, nil
}


func (s *Store) GetAllUserStats() (map[string]int, error) {
	stats := make(map[string]int)

	keys, err := s.client.Keys(s.ctx, "ratelimit:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		userID := key[len("ratelimit:"):]
		count, err := s.client.Get(s.ctx, key).Int64()
		if err != nil {
			continue
		}
		stats[userID] = int(count)
	}

	return stats, nil
}
