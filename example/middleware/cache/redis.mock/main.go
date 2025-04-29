// Example of cache middleware with Redis storage in Quick
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cache"
)

// MockRedisClient is a simple mock implementation of the RedisClient interface
// In a real qlication, you would use a real Redis client like github.com/go-redis/redis/v8
type MockRedisClient struct {
	storage map[string]string
	expiry  map[string]time.Time
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		storage: make(map[string]string),
		expiry:  make(map[string]time.Time),
	}
}

// Get retrieves a value from the mock Redis storage
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if expiry, ok := m.expiry[key]; ok {
		if time.Now().After(expiry) {
			delete(m.storage, key)
			delete(m.expiry, key)
			return "", fmt.Errorf("key expired")
		}
	}

	if value, ok := m.storage[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("key not found")
}

// Set stores a value in the mock Redis storage
func (m *MockRedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.storage[key] = value
	if expiration > 0 {
		m.expiry[key] = time.Now().Add(expiration)
	}
	return nil
}

// Del removes a value from the mock Redis storage
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	var count int64
	for _, key := range keys {
		if _, ok := m.storage[key]; ok {
			delete(m.storage, key)
			delete(m.expiry, key)
			count++
		}
	}
	return count, nil
}

func main() {
	// Create a new Quick q
	q := quick.New()

	// Create a mock Redis client
	// In a real qlication, you would use a real Redis client
	redisClient := NewMockRedisClient()

	// Create Redis storage
	redisStorage, err := cache.NewRedisStorage(cache.RedisConfig{
		Client:    redisClient,
		TTL:       1 * time.Minute,
		KeyPrefix: "quick:cache",
	})
	if err != nil {
		panic(err)
	}

	// Use the cache middleware with Redis storage
	q.Use(cache.New(cache.Config{
		Expiration:           1 * time.Minute,
		CacheHeader:          "X-Cache-Status",
		StoreResponseHeaders: true,
		Storage:              redisStorage,
	}))

	// Route 1: Returns the current time
	q.Get("/time", func(c *quick.Ctx) error {
		return c.SendString("Current time: " + time.Now().Format(time.RFC1123))
	})

	// Route 2: Returns a random number
	q.Get("/random", func(c *quick.Ctx) error {
		return c.SendString("Random value: " + time.Now().Format("15:04:05.000"))
	})

	// Route 3: Returns JSON data
	q.Get("/profile", func(c *quick.Ctx) error {
		return c.JSON(quick.M{
			"user":  "jeffotoni",
			"since": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// Start the server
	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("Using Redis storage (mock implementation for this example)")
	fmt.Println("Try these endpoints:")
	fmt.Println("  - GET /time (cached for 1 minute)")
	fmt.Println("  - GET /random (cached for 1 minute)")
	fmt.Println("  - GET /profile (cached for 1 minute)")
	fmt.Println("Check the X-Cache-Status header in the response to see if it's a HIT or MISS")
	q.Listen(":3000")
}
