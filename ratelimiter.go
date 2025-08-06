package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	clients map[string]*ClientBucket // Maps IP addresses to their buckets
	mutex   sync.RWMutex             // Protects the map from concurrent access
}

type ClientBucket struct {
	tokens     int       // Current number of tokens (requests) available
	lastRefill time.Time // Last time the bucket was refilled
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientBucket),
	}

	// Start the cleanup goroutine
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	// Create a timer that fires every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // Clean up the ticker when function exits

	// Run forever (until the program ends)
	for range ticker.C { // Runs every 5 minutes
		rl.mutex.Lock() // Lock the map for exclusive access
		now := time.Now()

		// Check each client bucket
		for ip, bucket := range rl.clients {
			// If bucket hasn't been touched in 10 minutes, delete it
			if now.Sub(bucket.lastRefill) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}

		rl.mutex.Unlock() // Release the lock
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()         // Lock the map for exclusive access
	defer rl.mutex.Unlock() // Ensure we unlock when done

	now := time.Now()
	bucket, exists := rl.clients[ip]

	// If the client doesn't have a bucket, create one
	if !exists {
		bucket = &ClientBucket{
			tokens:     10, // Start with 10 requests allowed
			lastRefill: now,
		}
		rl.clients[ip] = bucket // Add to filing cabinet
	}

	// Refill tokens (1 token per 6 seconds = 10 per minute)
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed.Seconds() / 6) // determine how many tokens to add

	if tokensToAdd > 0 { // add the tokens if any
		// Ensure we don't exceed the maximum of 10 tokens
		bucket.tokens = min(bucket.tokens+tokensToAdd, 10) // max 10 tokens
		bucket.lastRefill = now
	}

	if bucket.tokens > 0 { // If we have tokens left
		bucket.tokens-- // Take away 1 token
		return true     // Allow the request
	}

	return false // No tokens left, deny the request
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// min function for Go versions that don't have it built-in
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
