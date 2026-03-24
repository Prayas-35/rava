package middleware

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TokenBucket struct {
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

type RateLimiter struct {
	buckets        map[string]*TokenBucket
	mu             sync.RWMutex
	capacity       float64
	refillRate     float64
	refillInterval time.Duration
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		buckets:        make(map[string]*TokenBucket),
		capacity:       10.0,
		refillRate:     10.0,
		refillInterval: time.Second,
	}

	go rl.cleanupOldBuckets()
	return rl
}

func (rl *RateLimiter) hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return fmt.Sprintf("%x", hash)
}

func (rl *RateLimiter) allowToken(hashedKey string) bool {
	rl.mu.Lock()
	bucket, exists := rl.buckets[hashedKey]
	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.capacity,
			lastRefill: time.Now(),
		}
		rl.buckets[hashedKey] = bucket
	}
	rl.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := elapsed.Seconds() * rl.refillRate
	bucket.tokens = min(bucket.tokens+tokensToAdd, rl.capacity)
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	return false
}

func (rl *RateLimiter) cleanupOldBuckets() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			bucket.mu.Lock()
			if now.Sub(bucket.lastRefill) > 10*time.Minute {
				delete(rl.buckets, key)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

var globalRateLimiter = NewRateLimiter()

func RateLimiting() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Next()
		}

		projectID, err := service.ResolveProjectIDByAPIKey(c.Context(), apiKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or revoked API key")
		}

		hashedKey := globalRateLimiter.hashAPIKey(apiKey)

		if !globalRateLimiter.allowToken(hashedKey) {
			c.Set("Retry-After", "0.1")
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded: 10 requests per second")
		}

		_ = projectID

		return c.Next()
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
