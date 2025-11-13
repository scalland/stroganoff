package ratelimit

import (
	"sync"
	"time"

	"github.com/yourusername/stroganoff/internal/config"
)

// Limiter implements token bucket algorithm for rate limiting
type Limiter struct {
	mu       sync.RWMutex
	buckets  map[string]*bucket
	ticker   *time.Ticker
	stopCh   chan struct{}
}

type bucket struct {
	tokens    float64
	lastReset time.Time
}

// NewLimiter creates a new rate limiter
func NewLimiter() *Limiter {
	limiter := &Limiter{
		buckets: make(map[string]*bucket),
		stopCh:  make(chan struct{}),
	}

	// Start cleanup goroutine
	limiter.startCleanup()
	return limiter
}

// Allow checks if a request from the given identifier is allowed
func (l *Limiter) Allow(identifier string) bool {
	cfg := config.GetInstance().GetAPI()

	if cfg.RateLimit <= 0 {
		return true // Rate limiting disabled
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[identifier]
	now := time.Now()

	if !exists {
		b = &bucket{
			tokens:    float64(cfg.RateLimit),
			lastReset: now,
		}
		l.buckets[identifier] = b
	}

	// Refill bucket based on elapsed time
	window := time.Duration(cfg.RateLimitWindow) * time.Second
	elapsed := now.Sub(b.lastReset)
	refillRate := float64(cfg.RateLimit) / window.Seconds()
	b.tokens += refillRate * elapsed.Seconds()

	if b.tokens > float64(cfg.RateLimit) {
		b.tokens = float64(cfg.RateLimit)
	}

	b.lastReset = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

// startCleanup periodically cleans up old buckets
func (l *Limiter) startCleanup() {
	l.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-l.stopCh:
				l.ticker.Stop()
				return
			case <-l.ticker.C:
				l.cleanup()
			}
		}
	}()
}

func (l *Limiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	maxAge := 1 * time.Hour

	for identifier, b := range l.buckets {
		if now.Sub(b.lastReset) > maxAge {
			delete(l.buckets, identifier)
		}
	}
}

// Stop stops the rate limiter
func (l *Limiter) Stop() {
	close(l.stopCh)
}

// Reset resets the bucket for an identifier
func (l *Limiter) Reset(identifier string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, identifier)
}
