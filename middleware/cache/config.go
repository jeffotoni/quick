package cache

import (
	"time"

	"github.com/jeffotoni/quick"
)

// Storage defines the interface for cache storage implementations.
type Storage interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// Config defines the configuration options for the cache middleware.
type Config struct {
	// Expiration is the default duration after which cached items will expire.
	// Default is 1 minute.
	Expiration time.Duration

	// ExpirationGenerator is a function that returns a custom TTL for each request.
	// If provided, it takes precedence over the Expiration setting.
	ExpirationGenerator func(c *quick.Ctx, cfg *Config) time.Duration

	// KeyGenerator is a function that generates a unique cache key for each request.
	// Default is to use the request path.
	KeyGenerator func(c *quick.Ctx) string

	// CacheHeader is the name of the header that indicates cache status.
	// Default is "X-Cache-Status".
	CacheHeader string

	// CacheControl enables respecting Cache-Control headers from clients.
	// Default is true.
	CacheControl bool

	// StoreResponseHeaders determines whether to cache and restore response headers.
	// Default is true.
	StoreResponseHeaders bool

	// MaxBytes is the maximum size in bytes for a response to be cached.
	// Default is 1MB.
	MaxBytes int

	// Methods is a list of HTTP methods to cache.
	// Default is GET and HEAD.
	Methods []string

	// CacheInvalidator is a function that determines whether to skip the cache
	// for a specific request, effectively invalidating it.
	CacheInvalidator func(c *quick.Ctx) bool

	// Next is a function that determines whether to skip the middleware.
	Next func(c *quick.Ctx) bool

	// Storage is the cache storage engine to use.
	Storage Storage

	// OnHit is triggered when a cache hit occurs (optional).
	OnHit func(key string)

	// OnMiss is triggered when a cache miss occurs (optional).
	OnMiss func(key string)

	// OnCacheHit is called when a cache entry is successfully served.
	OnCacheHit func(c *quick.Ctx, key string)

	// OnCacheSet is called when a response is successfully cached.
	OnCacheSet func(c *quick.Ctx, key string)
}

// defaultConfig returns the default configuration for the cache middleware.
var defaultConfig = Config{
	Expiration:           1 * time.Minute,
	ExpirationGenerator:  nil,
	KeyGenerator:         nil,
	CacheHeader:          "X-Cache-Status",
	CacheControl:         true,
	StoreResponseHeaders: true,
	MaxBytes:             1024 * 1024, // 1MB
	Methods:              []string{quick.MethodGet, quick.MethodHead},
	CacheInvalidator:     nil,
	Next:                 nil,
	OnHit:                nil,
	OnMiss:               nil,
}

// cacheEntry represents a cached HTTP response.
type cacheEntry struct {
	Body         []byte
	StatusCode   int
	Headers      map[string][]string
	ContentType  string
	Expiration   time.Time
	LastAccessed time.Time
	CreatedAt    time.Time
}
