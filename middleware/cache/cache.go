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
			populateQueryParams(c)

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
			if shouldInvalidateCache(c, &cfg, key) {
				return next.ServeQuick(c)
			}

			// Check if client sent Cache-Control: no-cache
			if shouldBypassCache(c, &cfg) {
				return next.ServeQuick(c)
			}

			// Try to get from cache
			if cached, found := cfg.Storage.Get(key); found {
				if cfg.OnHit != nil {
					cfg.OnHit(key)
				}
				if cfg.OnCacheHit != nil {
					cfg.OnCacheHit(c, key)
				}
				return handleCacheHit(c, &cfg, cached)
			}
			if cfg.OnMiss != nil {
				cfg.OnMiss(key)
			}

			// Handle cache miss
			return handleCacheMiss(c, &cfg, next, key)
		})
	}
}

// populateQueryParams ensures query parameters are populated from the URL
func populateQueryParams(c *quick.Ctx) {
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
}

// shouldInvalidateCache checks if the cache should be invalidated
func shouldInvalidateCache(c *quick.Ctx, cfg *Config, key string) bool {
	if cfg.CacheInvalidator != nil && cfg.CacheInvalidator(c) {
		// Delete the cache entry
		cfg.Storage.Delete(key)
		c.Set(cfg.CacheHeader, "INVALIDATED")
		return true
	}
	return false
}

// shouldBypassCache checks if the cache should be bypassed
func shouldBypassCache(c *quick.Ctx, cfg *Config) bool {
	if cfg.CacheControl && c.Get("Cache-Control") == "no-cache" {
		c.Set(cfg.CacheHeader, "BYPASS")
		return true
	}
	return false
}

// handleCacheHit processes a cache hit
func handleCacheHit(c *quick.Ctx, cfg *Config, cached interface{}) error {
	entry := cached.(*cacheEntry)

	// Check if the entry is expired
	if time.Now().After(entry.Expiration) {
		cfg.Storage.Delete(cfg.KeyGenerator(c))
		c.Set(cfg.CacheHeader, "EXPIRED")
		return nil
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

// handleCacheMiss processes a cache miss
func handleCacheMiss(c *quick.Ctx, cfg *Config, next quick.Handler, key string) error {
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

	// If there was an error or the response is too large, handle without caching
	if err != nil || responseWriter.buffer.Len() > cfg.MaxBytes {
		return handleResponseWithoutCaching(c, cfg, responseWriter, err)
	}

	// Process and cache the successful response
	return processAndCacheResponse(c, cfg, responseWriter, key)
}

// handleResponseWithoutCaching handles responses that shouldn't be cached
func handleResponseWithoutCaching(c *quick.Ctx, cfg *Config, responseWriter *responseCapture, err error) error {
	// If there was no error, write the response
	if err == nil {
		// Copy headers to the original response
		copyHeaders(c.Response, responseWriter.headers)

		// Cache the response if it's not too large and there was no error
		if responseWriter.buffer.Len() <= cfg.MaxBytes {
			cacheResponse(c, cfg, responseWriter)
		}

		// Write the body directly
		_, err = c.Response.Write(responseWriter.buffer.Bytes())
	}
	return err
}

// processAndCacheResponse processes and caches a successful response
func processAndCacheResponse(c *quick.Ctx, cfg *Config, responseWriter *responseCapture, key string) error {
	// Determine expiration time
	expiration := calculateExpiration(c, cfg)

	// Create a new cache entry
	entry := buildCacheEntry(c, responseWriter, cfg, expiration)

	// Store response headers if configured
	if cfg.StoreResponseHeaders {
		entry.Headers = responseWriter.headers
	}

	// Store in cache
	cfg.Storage.Set(key, entry, time.Until(expiration))

	if cfg.OnCacheSet != nil {
		cfg.OnCacheSet(c, key)
	}

	// Copy important headers to avoid WriteHeader conflicts
	copyImportantHeaders(c.Response, responseWriter.headers)

	// Write the body directly and return
	_, err := c.Response.Write(responseWriter.buffer.Bytes())
	return err
}

// calculateExpiration determines the expiration time for a cache entry
func calculateExpiration(c *quick.Ctx, cfg *Config) time.Time {
	if cfg.ExpirationGenerator != nil {
		return time.Now().Add(cfg.ExpirationGenerator(c, cfg))
	}
	return time.Now().Add(cfg.Expiration)
}

// cacheResponse caches a response
func cacheResponse(c *quick.Ctx, cfg *Config, responseWriter *responseCapture) {
	expiration := calculateExpiration(c, cfg)
	entry := buildCacheEntry(c, responseWriter, cfg, expiration)

	// Store response headers if configured
	if cfg.StoreResponseHeaders {
		entry.Headers = responseWriter.headers
	}

	// Store in cache
	cfg.Storage.Set(cfg.KeyGenerator(c), entry, time.Until(expiration))
}

// copyHeaders copies all headers from src to dst
func copyHeaders(dst http.ResponseWriter, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Header().Set(key, value)
		}
	}
}

// copyImportantHeaders copies only important headers (Content-Type, Content-Length, X-*)
func copyImportantHeaders(dst http.ResponseWriter, src http.Header) {
	for key, values := range src {
		if key == "Content-Type" || key == "Content-Length" || strings.HasPrefix(key, "X-") {
			for _, value := range values {
				dst.Header().Set(key, value)
			}
		}
	}
}

func buildCacheEntry(c *quick.Ctx, w *responseCapture, cfg *Config, exp time.Time) *cacheEntry {
	contentType := getHeader(w.headers, "Content-Type")
	if contentType == "" {
		contentType = w.ResponseWriter.Header().Get("Content-Type")
	}

	if contentType == "" {
		// last fallback: detect from body or use text/plain
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

// Flush implements http.Flusher interface for SSE support
func (r *responseCapture) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
