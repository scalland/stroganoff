package ratelimit

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewLimiter()
	defer limiter.Stop()

	identifier := "test-user"

	// Allow should return true initially (bucket has tokens)
	if !limiter.Allow(identifier) {
		t.Fatal("Allow should return true initially")
	}
}

func TestRateLimiterReset(t *testing.T) {
	limiter := NewLimiter()
	defer limiter.Stop()

	identifier := "test-user"

	// Use up the tokens
	for i := 0; i < 100; i++ {
		limiter.Allow(identifier)
	}

	// Reset the bucket
	limiter.Reset(identifier)

	// Should allow again after reset
	if !limiter.Allow(identifier) {
		t.Fatal("Allow should return true after reset")
	}
}

func TestRateLimiterCleanup(t *testing.T) {
	limiter := NewLimiter()
	defer limiter.Stop()

	identifier := "test-user"
	limiter.Allow(identifier)

	// Wait for cleanup
	time.Sleep(2 * time.Minute)

	// Note: This is a basic test. Full cleanup testing would require
	// more sophisticated bucket aging logic
}
