// Package cache provides a middleware for the Quick web framework
// that implements an in-memory caching system for HTTP responses.
package cache

import (
	"sync"
	"time"
)

const (
	// DefaultExpiration uses the cache's default TTL if no specific expiration is set.
	DefaultExpiration time.Duration = 0

	// NoExpiration indicates that the cached item should never expire.
	NoExpiration time.Duration = -1

	// numShards defines the number of cache partitions to allow concurrent access.
	numShards = 8

	// ringSize sets the size of the expiration ring buffer for tracking expired items.
	ringSize = 4096
)

// ringNode represents an entry in the expiration ring buffer.
type ringNode struct {
	key     uint32 // Hashed key
	expires int64  // Expiration timestamp in nanoseconds
}

// shard is a partition of the cache with its own locking mechanism.
type shard struct {
	mu       sync.RWMutex     // Mutex for concurrent access
	items    map[uint32]*Item // Cached items
	ringBuf  []ringNode       // Ring buffer for tracking expiration
	ringHead int              // Current position in the ring buffer
}

// Item represents a single cache entry.
type Item struct {
	value   interface{} // Stored value
	expires int64       // Expiration timestamp
}

// Cache is a sharded in-memory cache with expiration handling.
type Cache struct {
	shards [numShards]*shard // Array of shards to reduce contention
	ttl    time.Duration     // Default time-to-live for cache entries
}

// NewCache creates a new instance of Cache with the specified default TTL.
// If the TTL is greater than 0, a cleanup goroutine is started to periodically remove expired items.
func NewCache(ttlStr ...time.Duration) *Cache {
	var ttl time.Duration
	if len(ttlStr) > 0 {
		// Use the first duration provided
		ttl = ttlStr[0]
	} else {
		// Fallback to DefaultExpiration if no parameter is passed
		ttl = DefaultExpiration
	}
	c := &Cache{ttl: ttl}
	for i := 0; i < numShards; i++ {
		c.shards[i] = &shard{
			items:   make(map[uint32]*Item),
			ringBuf: make([]ringNode, ringSize),
		}
	}
	if ttl > 0 {
		go c.cleanup()
	}
	return c
}

// hashKey computes a simple FNV-1a hash from the string key.
// The hash ensures even distribution across shards.
func (c *Cache) hashKey(key string) uint32 {
	var h uint32
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return h
}

// getShard selects the shard based on the hashed key.
// This helps in distributing load and reducing lock contention.
func (c *Cache) getShard(k uint32) *shard {
	return c.shards[k%numShards]
}

// Set inserts a value into the cache with an optional TTL.
// If `ttl` is set to `DefaultExpiration`, the cache's default TTL is applied.
// If `ttl` is set to `NoExpiration`, the item never expires.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	var exp int64
	if ttl == DefaultExpiration {
		ttl = c.ttl
	}
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	hashed := c.hashKey(key)
	sh := c.getShard(hashed)

	sh.mu.Lock()
	sh.items[hashed] = &Item{value: value, expires: exp}
	sh.ringBuf[sh.ringHead] = ringNode{key: hashed, expires: exp}
	sh.ringHead = (sh.ringHead + 1) % ringSize
	sh.mu.Unlock()
}

// Get retrieves a value from the cache.
// Returns the stored value and a boolean indicating if the key was found.
// If the item has expired, it is removed from the cache and (nil, false) is returned.
func (c *Cache) Get(key string) (interface{}, bool) {
	hashed := c.hashKey(key)
	sh := c.getShard(hashed)

	sh.mu.RLock()
	item, exists := sh.items[hashed]
	sh.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		c.Delete(key) // Remove expired item
		return nil, false
	}

	return item.value, true
}

// Delete removes an item from the cache.
// If the key does not exist, no action is taken.
func (c *Cache) Delete(key string) {
	hashed := c.hashKey(key)
	sh := c.getShard(hashed)

	sh.mu.Lock()
	delete(sh.items, hashed)
	sh.mu.Unlock()
}

// cleanup runs periodically to remove expired items from the cache.
// This function runs as a background goroutine and checks for expired items
// at intervals of `ttl / 2`, ensuring efficient memory management.
func (c *Cache) cleanup() {
	tick := time.NewTicker(c.ttl / 2)
	defer tick.Stop()

	for range tick.C {
		now := time.Now().UnixNano()
		for _, sh := range c.shards {
			sh.mu.Lock()
			for i := 0; i < ringSize; i++ {
				node := &sh.ringBuf[i]
				if node.expires > 0 && now > node.expires {
					delete(sh.items, node.key)
					node.expires = 0
				}
			}
			sh.mu.Unlock()
		}
	}
}
