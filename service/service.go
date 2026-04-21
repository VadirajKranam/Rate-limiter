package service

import (
	"log"

	"github.com/go-rate-limiter/db"
)

const (
	MaxRequests = 5
	WindowSize  = 60
)

type RateLimiter struct {
	store *db.Store
}

func New(store *db.Store) *RateLimiter {
	return &RateLimiter{store: store}
}

func (rl *RateLimiter) Allow(userID string) bool {
	count, err := rl.store.IncrementAndGet(userID, WindowSize)
	if err != nil {
		log.Printf("Error checking rate limit: %v", err)
		return false
	}

	return count <= MaxRequests
}

func (rl *RateLimiter) GetStats() map[string]int {
	stats, err := rl.store.GetAllUserStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		return make(map[string]int)
	}
	return stats
}
