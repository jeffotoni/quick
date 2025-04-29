// Package cache provides a middleware for the Quick web framework
// that implements an in-memory caching system for HTTP responses.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// RedisClient defines the interface for a Redis client.
// This allows for flexibility in the Redis client implementation.
//
// The RedisClient interface follows the dependency inversion principle,
// where the middleware depends on an abstraction (this interface) rather than
// a concrete implementation. This design provides several benefits:
//
// 1. Flexibility: Any Redis client library can be used as long as it's adapted
//    to implement this interface.
//
// 2. Decoupling: The middleware doesn't need to know the specific details of
//    how Redis is accessed, only that something implements these methods.
//
// 3. Testability: Makes it easy to create mocks for unit testing the middleware.
//
// 4. Adaptability: If the Redis API changes or if you want to switch to another
//    storage system, you only need to modify the adapter, not the middleware.
//
// 5. Dependency control: Avoids direct dependencies on external libraries,
//    which could cause version conflicts or unnecessarily increase package size.
//
// To use this with a specific Redis client library, create an adapter that
// implements this interface and bridges between the library and the middleware.
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
}

// RedisStorage is a cache storage implementation that uses Redis.
type RedisStorage struct {
	client      RedisClient
	ttl         time.Duration
	keyPrefix   string
	ctx         context.Context
}

// RedisConfig defines the configuration for Redis storage.
type RedisConfig struct {
	// Client is the Redis client implementation to use.
	Client RedisClient

	// TTL is the default time-to-live for cache entries.
	TTL time.Duration

	// KeyPrefix is an optional prefix for all cache keys.
	KeyPrefix string
}

// NewRedisStorage creates a new Redis storage instance.
func NewRedisStorage(config RedisConfig) (*RedisStorage, error) {
	if config.Client == nil {
		return nil, errors.New("redis client is required")
	}

	if config.TTL <= 0 {
		config.TTL = 5 * time.Minute // Default TTL
	}

	return &RedisStorage{
		client:    config.Client,
		ttl:       config.TTL,
		keyPrefix: config.KeyPrefix,
		ctx:       context.Background(),
	}, nil
}

// Set stores a value in Redis with the specified TTL.
func (r *RedisStorage) Set(key string, value interface{}, ttl time.Duration) {
	if ttl == DefaultExpiration {
		ttl = r.ttl
	}

	// Prefix the key if a prefix is configured
	if r.keyPrefix != "" {
		key = r.keyPrefix + ":" + key
	}

	// Marshal the value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	// Store in Redis
	_ = r.client.Set(r.ctx, key, string(data), ttl)
}

// Get retrieves a value from Redis.
func (r *RedisStorage) Get(key string) (interface{}, bool) {
	// Prefix the key if a prefix is configured
	if r.keyPrefix != "" {
		key = r.keyPrefix + ":" + key
	}

	// Get from Redis
	data, err := r.client.Get(r.ctx, key)
	if err != nil {
		return nil, false
	}

	// Unmarshal the JSON data
	var entry cacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, false
	}

	return &entry, true
}

// Delete removes a value from Redis.
func (r *RedisStorage) Delete(key string) {
	// Prefix the key if a prefix is configured
	if r.keyPrefix != "" {
		key = r.keyPrefix + ":" + key
	}

	// Delete from Redis
	_, _ = r.client.Del(r.ctx, key)
}
