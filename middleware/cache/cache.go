// Package cache provides a middleware for the Quick web framework
// that implements an in-memory caching system for HTTP responses.
//
// This middleware intercepts responses and stores them in memory,
// serving subsequent identical requests directly from the cache
// without executing handlers. It supports configurable TTL,
// custom key generation, conditional invalidation, and more.
//
// The cache implementation is based on a high-performance algorithm
// copied from the gocache library, optimized for minimal lock contention
// and efficient memory usage.
package cache

import (
	"bytes"
	"net/http"
	"strings"
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
	// Default is an in-memory cache.
	Storage Storage
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
}

// New creates a new cache middleware with the provided configuration.
//
// Example usage:
//
//	app.Use(cache.New(cache.Config{
//		Expiration: 5 * time.Minute,
//		KeyGenerator: func(c *quick.Ctx) string {
//			return c.Path() + "?user=" + c.Query("user")
//		},
//	}))
func New(config ...Config) func(next quick.Handler) quick.Handler {
	// Apply default config
	cfg := defaultConfig
	if len(config) > 0 {
		cfg = config[0]

		// Apply defaults for any zero values
		if cfg.Expiration <= 0 {
			cfg.Expiration = defaultConfig.Expiration
		}
		if cfg.CacheHeader == "" {
			cfg.CacheHeader = defaultConfig.CacheHeader
		}
		if cfg.MaxBytes <= 0 {
			cfg.MaxBytes = defaultConfig.MaxBytes
		}
		if len(cfg.Methods) == 0 {
			cfg.Methods = defaultConfig.Methods
		}
	}

	// Initialize default key generator if not provided
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = func(c *quick.Ctx) string {
			// Use only the path as the default key
			return c.Path()
		}
	}

	// Initialize default Next function if not provided
	if cfg.Next == nil {
		cfg.Next = func(c *quick.Ctx) bool {
			return false
		}
	}

	// Initialize storage if not provided
	if cfg.Storage == nil {
		cfg.Storage = NewCache(cfg.Expiration)
	}

	// Create a slice of allowed methods for faster lookup
	methodMap := make(map[string]struct{}, len(cfg.Methods))
	for _, method := range cfg.Methods {
		methodMap[method] = struct{}{}
	}

	// Return the middleware handler
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Ensure query parameters are populated from URL
			if c.Query == nil {
				c.Query = make(map[string]string)
			}
			// Extract query parameters from URL
			query := c.Request.URL.Query()
			for k, v := range query {
				if len(v) > 0 {
					c.Query[k] = v[0]
				}
			}

			// Skip middleware if Next returns true
			if cfg.Next != nil && cfg.Next(c) {
				return next.ServeQuick(c)
			}

			// Only cache specified methods
			method := c.Method()
			if _, exists := methodMap[method]; !exists {
				return next.ServeQuick(c)
			}

			// Generate cache key
			key := cfg.KeyGenerator(c)

			// Check if cache should be invalidated
			if cfg.CacheInvalidator != nil && cfg.CacheInvalidator(c) {
				// Delete the cache entry
				cfg.Storage.Delete(key)
				c.Set(cfg.CacheHeader, "INVALIDATED")
				// Process the request and don't cache the response
				return next.ServeQuick(c)
			}

			// Check if client sent Cache-Control: no-cache
			if cfg.CacheControl && c.Get("Cache-Control") == "no-cache" {
				c.Set(cfg.CacheHeader, "BYPASS")
				return next.ServeQuick(c)
			}

			// Try to get from cache
			if cached, found := cfg.Storage.Get(key); found {
				entry := cached.(*cacheEntry)

				// Check if the entry is expired
				if time.Now().After(entry.Expiration) {
					cfg.Storage.Delete(key)
					c.Set(cfg.CacheHeader, "EXPIRED")
					// Process the request and don't cache the response
					return next.ServeQuick(c)
				}

				// Set the X-Cache header to indicate a cache hit
				c.Set(cfg.CacheHeader, "HIT")
				c.Set("X-Cache-Source", "memory")
				c.Set("X-Cache-Expires-At", entry.Expiration.Format(time.RFC3339))

				// Set headers directly on the response to avoid WriteHeader calls
				// Set the content type and other headers if configured
				if cfg.StoreResponseHeaders {
					for key, values := range entry.Headers {
						for _, value := range values {
							c.Response.Header().Set(key, value)
						}
					}
				} else if entry.ContentType != "" {
					c.Response.Header().Set("Content-Type", entry.ContentType)
				}

				// Set the status code directly on the response writer
				c.Response.WriteHeader(entry.StatusCode)

				// Write the cached response body directly to avoid framework methods
				_, err := c.Response.Write(entry.Body)
				return err
			}

			// Set X-Cache header to indicate a cache miss
			c.Set(cfg.CacheHeader, "MISS")

			// Create a response capture to store the response
			responseWriter := &responseCapture{
				ResponseWriter: c.Response,
				buffer:         bytes.NewBuffer(nil),
				headers:        make(http.Header),
			}

			// Replace the original response writer with our capturing one
			originalWriter := c.Response
			c.Response = responseWriter

			// Process the request with the next handler
			err := next.ServeQuick(c)

			// Restore the original response writer
			c.Response = originalWriter

			// If there was an error or the response is too large, don't cache
			if err != nil || responseWriter.buffer.Len() > cfg.MaxBytes {
				if err == nil {
					_, err = c.Response.Write(responseWriter.buffer.Bytes())
				}
				// Write the captured response to the original writer without calling WriteHeader again
				// Set headers directly on the response
				for key, values := range responseWriter.headers {
					for _, value := range values {
						c.Response.Header().Set(key, value)
					}
				}

				// Cache the response if it's not too large
				if responseWriter.buffer.Len() <= cfg.MaxBytes {
					// Determine expiration time
					var expiration time.Time
					if cfg.ExpirationGenerator != nil {
						expiration = time.Now().Add(cfg.ExpirationGenerator(c, &cfg))
					} else {
						expiration = time.Now().Add(cfg.Expiration)
					}

					// Create cache entry
					// entry := &cacheEntry{
					// 	Body:        responseWriter.buffer.Bytes(),
					// 	StatusCode:  responseWriter.statusCode,
					// 	Expiration:  expiration,
					// 	ContentType: responseWriter.Header().Get("Content-Type"),
					// }

					entry := buildCacheEntry(c, responseWriter, &cfg, expiration)

					// Store response headers if configured
					if cfg.StoreResponseHeaders {
						entry.Headers = responseWriter.headers
					}

					// Store in cache
					cfg.Storage.Set(key, entry, time.Until(expiration))
				}

				// Write the body directly
				_, err = c.Response.Write(responseWriter.buffer.Bytes())
				return err
			}

			// Determine expiration time
			var expiration time.Time
			if cfg.ExpirationGenerator != nil {
				expiration = time.Now().Add(cfg.ExpirationGenerator(c, &cfg))
			} else {
				expiration = time.Now().Add(cfg.Expiration)
			}

			// Create a new cache entry
			entry := buildCacheEntry(c, responseWriter, &cfg, expiration)

			// Store response headers if configured
			if cfg.StoreResponseHeaders {
				entry.Headers = responseWriter.headers
			}

			// Store in cache
			cfg.Storage.Set(key, entry, time.Until(expiration))

			// Completely bypass the framework's status setting to avoid WriteHeader conflicts
			// Just write the response body directly to the underlying http.ResponseWriter

			// Copy any important headers
			for key, values := range responseWriter.headers {
				if key == "Content-Type" || key == "Content-Length" || strings.HasPrefix(key, "X-") {
					for _, value := range values {
						c.Response.Header().Set(key, value)
					}
				}
			}

			// Write the body directly and return
			_, err = c.Response.Write(responseWriter.buffer.Bytes())
			return err
		})
	}
}

func buildCacheEntry(c *quick.Ctx, w *responseCapture, cfg *Config, exp time.Time) *cacheEntry {
	contentType := getHeader(w.headers, "Content-Type")
	if contentType == "" {
		contentType = w.ResponseWriter.Header().Get("Content-Type")
	}

	if contentType == "" {
		// último fallback: detecta a partir do corpo ou usa text/plain
		if b := w.buffer.Bytes(); len(b) > 0 {
			contentType = http.DetectContentType(b)
		} else {
			contentType = "text/plain; charset=utf-8"
		}
	}

	return &cacheEntry{
		Body:         w.buffer.Bytes(),
		StatusCode:   w.statusCode,
		Expiration:   exp,
		ContentType:  contentType,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}
}

func getHeader(h map[string][]string, key string) string {
	if values, ok := h[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

// responseCapture is a wrapper around http.ResponseWriter that captures
// the response for caching.
type responseCapture struct {
	http.ResponseWriter
	statusCode    int
	buffer        *bytes.Buffer
	headers       http.Header
	headerWritten bool // Flag to track if WriteHeader has been called
}

// WriteHeader captures the status code.
func (r *responseCapture) WriteHeader(statusCode int) {
	if r.headerWritten {
		return
	}
	r.statusCode = statusCode
	r.headerWritten = true

	for k, vv := range r.headers {
		for _, v := range vv {
			r.ResponseWriter.Header().Add(k, v)
		}
	}
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body.
func (r *responseCapture) Write(b []byte) (int, error) {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
	r.buffer.Write(b)
	return r.ResponseWriter.Write(b)
}

// Header captures response headers.
func (r *responseCapture) Header() http.Header {
	return r.ResponseWriter.Header()
}

// Set adds a header to the response.
func (r *responseCapture) Set(key, value string) {
	r.headers[key] = append(r.headers[key], value)
	r.ResponseWriter.Header().Set(key, value)
}

// Add appends a header to the response.
func (r *responseCapture) Add(key, value string) {
	r.headers[key] = append(r.headers[key], value)
	r.ResponseWriter.Header().Add(key, value)
}

// Del removes a header from the response.
func (r *responseCapture) Del(key string) {
	delete(r.headers, key)
	r.ResponseWriter.Header().Del(key)
}
