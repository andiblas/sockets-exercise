package ratelimiter

import (
	"net"
	"sync"
	"time"
)

type TokenBucket struct {
	mu          sync.Mutex
	rate        int           // Number of requests allowed per time period
	per         time.Duration // Time period
	lastRequest map[string]time.Time
}

func NewTokenBucket(rate int, per time.Duration) *TokenBucket {
	return &TokenBucket{
		mu:          sync.Mutex{},
		rate:        rate,
		per:         per,
		lastRequest: make(map[string]time.Time),
	}
}

func (rl *TokenBucket) Allow(conn net.Conn) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	clientAddr := conn.RemoteAddr().String()

	lastReqTime, ok := rl.lastRequest[clientAddr]
	if !ok || time.Since(lastReqTime) > rl.per {
		rl.lastRequest[clientAddr] = time.Now()
		return true
	}
	return false
}
