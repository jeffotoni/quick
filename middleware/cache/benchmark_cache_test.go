package cache

import (
	"fmt"
	"testing"
	"time"
	"github.com/jeffotoni/quick"
)

func BenchmarkCacheMiddleware_Cache_Hit(b *testing.B) {
	q := quick.New()
	q.Use(New())

	// Simple handler that returns a constant response
	q.Get("/bench", func(c *quick.Ctx) error {
		return c.String("hit")
	})

	// First request to populate the cache
	_, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/bench",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/bench",
		})
		if err != nil {
			b.Fatalf("Cache hit error: %v", err)
		}
		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			b.Fatalf("Expected HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
	}
}

func BenchmarkCacheMiddleware_Cache_Miss(b *testing.B) {
	q := quick.New()
	q.Use(New(Config{
		KeyGenerator: func(c *quick.Ctx) string {
			// Ensures a unique cache key for every request (always a MISS)
			return c.Path() + "?" + c.Request.URL.RawQuery
		},
	}))

	q.Get("/benchmiss", func(c *quick.Ctx) error {
		return c.String("miss")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uri := "/benchmiss?i=" + fmt.Sprint(i) // Unique URI for each request
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    uri,
		})
		if err != nil {
			b.Fatalf("Cache miss error: %v", err)
		}
		if resp.Response().Header.Get("X-Cache-Status") != "MISS" {
			b.Fatalf("Expected MISS, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
	}
}

// Mock that simulates RedisStorage implementing the Storage interface
// (does not depend on context or real Redis)
type MockRedisStorage struct {
	store map[string]interface{}
}

func NewMockRedisStorage() *MockRedisStorage {
	return &MockRedisStorage{store: make(map[string]interface{})}
}

func (m *MockRedisStorage) Get(key string) (interface{}, bool) {
	val, ok := m.store[key]
	return val, ok
}

func (m *MockRedisStorage) Set(key string, value interface{}, ttl time.Duration) {
	m.store[key] = value
}

func (m *MockRedisStorage) Delete(key string) {
	delete(m.store, key)
}

func BenchmarkCacheMiddleware_Redis_Hit(b *testing.B) {
	mockRedis := NewMockRedisStorage()
	cfg := Config{
		Storage: mockRedis,
	}
	q := quick.New()
	q.Use(New(cfg))

	q.Get("/redisbench", func(c *quick.Ctx) error {
		return c.String("redis hit")
	})

	// First request to populate the mock Redis cache
	_, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/redisbench",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/redisbench",
		})
		if err != nil {
			b.Fatalf("Redis hit error: %v", err)
		}
		if resp.Response().Header.Get("X-Cache-Status") != "HIT" {
			b.Fatalf("Expected HIT, got %s", resp.Response().Header.Get("X-Cache-Status"))
		}
	}
}
