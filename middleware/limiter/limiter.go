// Package limiter provides middleware for rate limiting HTTP requests in Quick.
//
// This middleware controls the number of requests a client can make within a specified time window,
// helping to prevent abuse, protect APIs from excessive traffic, and improve overall system stability.
//
// Features:
// - Configurable maximum requests per time window.
// - Customizable key generator (e.g., per-IP, per-user, etc.).
// - Flexible expiration time for rate-limited requests.
// - Custom handler when the request limit is exceeded.
// - Uses sharded maps for efficient concurrency handling.
// - Periodic cleanup of expired request records to optimize memory usage.
package limiter

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"sync"
	"time"

	"github.com/jeffotoni/quick"
)

// Config defines the rate limiting configuration.
//
// Max: maximum number of requests allowed within the expiration window.
// Expiration: duration before the request counter resets.
// KeyGenerator: function to generate a unique key for each client (e.g., IP-based).
// LimitReached: function called when the client exceeds the request limit.
type Config struct {
	Max          int                     // Maximum requests allowed in the time window
	Expiration   time.Duration           // Time window for rate limiting
	KeyGenerator func(*quick.Ctx) string // Function to generate a unique key per client
	LimitReached func(*quick.Ctx) error  // Function executed when rate limit is exceeded
}

// client tracks individual request data for rate limiting.
type client struct {
	mu       sync.Mutex
	requests int       // Number of requests made in the current window
	expires  time.Time // When the current window expires
}

// RateLimiter manages all rate limiting logic, storing request counters across multiple shards.
type RateLimiter struct {
	config     Config      // Rate limiting settings
	shards     []*sync.Map // Sharded maps for distributing load and reducing contention
	shardCount uint32      // Number of shards for concurrency
}

// New creates a middleware constructor that returns a standard http.Handler wrapper.
//
// Usage:
//
//	q.Use(limiter.New(limiter.Config{
//	    Max:        3,
//	    Expiration: 2 * time.Second,
//	    KeyGenerator: func(c *quick.Ctx) string {
//	        // Return IP without port, or a fixed test key.
//	        return "testKey"
//	    },
//	    LimitReached: func(c *quick.Ctx) error {
//	        return c.Status(http.StatusTooManyRequests).SendString("Too many requests")
//	    },
//	}))
//
// The returned function integrates with Quick's middleware chain and enforces rate limits.
func New(config Config) func(http.Handler) http.Handler {
	rl := &RateLimiter{
		config:     config,
		shardCount: 256, // Default shard count for concurrency
	}

	// Initialize the shards
	rl.shards = make([]*sync.Map, rl.shardCount)
	for i := 0; i < int(rl.shardCount); i++ {
		rl.shards[i] = &sync.Map{}
	}

	// Start the background cleanup to remove expired clients
	go rl.startCleanup()

	// Return the middleware constructor
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a Quick context to allow KeyGenerator or LimitReached usage.
			c := &quick.Ctx{
				Response: w,
				Request:  r,
			}

			// Generate the key (e.g., IP) for this client.
			key := rl.config.KeyGenerator(c)
			//fmt.Println("[DEBUG] Key:", key)

			// Find the shard and load the client's request counter.
			shard := rl.getShard(key)
			now := time.Now()

			val, _ := shard.LoadOrStore(key, &client{
				requests: 0,
				expires:  now.Add(rl.config.Expiration),
			})
			cl := val.(*client)

			// Lock to safely update the request counter.
			cl.mu.Lock()
			defer cl.mu.Unlock()

			// Reset if the current window expired.
			if now.After(cl.expires) {
				//fmt.Println("[DEBUG] Resetting requests")
				cl.requests = 0
				cl.expires = now.Add(rl.config.Expiration)
			}

			// Increment the request count.
			cl.requests++
			//fmt.Printf("[DEBUG] requests=%d, max=%d\n", cl.requests, rl.config.Max)

			// If the client exceeded the limit, call LimitReached and stop.
			if cl.requests > rl.config.Max {
				//fmt.Println("[DEBUG] Limit reached: calling LimitReached")
				err := rl.config.LimitReached(c)
				if err != nil {
					fmt.Println("[DEBUG] LimitReached error:", err)
				}
				return
			}

			// Otherwise, pass to the next handler.
			next.ServeHTTP(w, r)
		})
	}
}

// getShard selects which shard map is used for the given key.
func (rl *RateLimiter) getShard(key string) *sync.Map {
	h := fnv.New32a()
	h.Write([]byte(key))
	return rl.shards[h.Sum32()%rl.shardCount]
}

// startCleanup periodically removes expired client entries.
func (rl *RateLimiter) startCleanup() {
	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()
	for range tick.C {
		rl.cleanup()
	}
}

// cleanup checks every shard and deletes clients whose expiration time has passed.
func (rl *RateLimiter) cleanup() {
	now := time.Now()
	for _, shard := range rl.shards {
		shard.Range(func(k, v interface{}) bool {
			cl := v.(*client)
			cl.mu.Lock()
			expired := now.After(cl.expires)
			cl.mu.Unlock()

			if expired {
				shard.Delete(k)
			}
			return true
		})
	}
}
