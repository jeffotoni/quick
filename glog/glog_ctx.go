// Package glog provides an optimized zero-allocation logger and contextual tracing helpers.
//
// This file extends glog's capabilities by providing context propagation utilities
// using a fluent API (`CtxBuilder`) that allows you to safely and efficiently inject
// metadata like Trace IDs into `context.Context`. This is useful for observability,
// distributed tracing, correlation IDs, or user/session propagation.
//
// Key goals:
//   - No heap allocations using string-based context keys and sync.Map caching
//   - Fluent builder for setting multiple fields and optional timeout
//   - Compatible with standard context propagation in Go HTTP servers, workers, etc.
//
// # Features:
//
//   - Zero-allocation context key management
//   - Fluent builder to set trace fields and timeout
//   - Safe fallback when context is nil
//   - Retrieve individual or all injected fields
//   - Optimized for performance and readability
//
// # Example:
//
//	import (
//		"fmt"
//		"time"
//		"github.com/jeffotoni/quick/glog"
//	)
//
//	func main() {
//		ctx, cancel := glog.CreateCtx().
//			Set("TraceID", "abc-123").
//			Set("UserID", "u42").
//			Timeout(5 * time.Second).
//			Build()
//		defer cancel()
//
//		trace := glog.GetCtx(ctx, "TraceID")
//		fmt.Println("Trace:", trace)
//
//		all := glog.GetCtxAll(ctx)
//		fmt.Println("All Fields:", all)
//	}
package glog

import (
	"context"
	"sync"
	"time"
)

// internalKeysKey is the reserved context key used to track which custom keys
// were injected by the context builder, enabling retrieval via GetCtxAll.
const internalKeysKey = "__glog_ctx_keys__"

// Predefined constants to reduce allocations and provide sane defaults.
const (
	defaultCtxTimeout = 30 * time.Second // Default timeout for built contexts
)

// contextKey is a private type to avoid key collisions in context.Context.
// A string-based key is used instead of a struct to reduce allocations.
type contextKey string

// Cached context keys using sync.Map to minimize allocations during key creation.
var (
	keyCache     sync.Map               // map[string]contextKey
	emptyContext = context.Background() // Reusable base context
)

// getCtxKey returns a cached context key or creates and stores a new one.
// This reduces allocations by reusing frequently accessed keys.
func getCtxKey(name string) contextKey {
	if name == "" {
		return contextKey("default")
	}

	if v, ok := keyCache.Load(name); ok {
		return v.(contextKey)
	}

	k := contextKey(name)
	keyCache.Store(name, k)
	return k
}

// CtxBuilder provides a fluent API for constructing a context.Context
// with multiple key-value fields and an optional timeout.
type CtxBuilder struct {
	fields  map[string]string
	timeout time.Duration
}

// CreateCtx initializes and returns a new context builder (CtxBuilder)
// with default timeout. Use .Set(...) and .Timeout(...) to add metadata.
//
// Example:
//
//	ctx, cancel := glog.CreateCtx().Set("TraceID", "abc").Build()
func CreateCtx() *CtxBuilder {
	return &CtxBuilder{
		fields:  make(map[string]string),
		timeout: defaultCtxTimeout,
	}
}

// Set injects a key-value string pair into the context.
// Keys and values must be non-empty.
func (b *CtxBuilder) Set(key, value string) *CtxBuilder {
	if key != "" && value != "" {
		b.fields[key] = value
	}
	return b
}

// Timeout sets a customized timeout duration for the context
func (b *CtxBuilder) Timeout(d time.Duration) *CtxBuilder {
	if d > 0 {
		b.timeout = d
	}
	return b
}

// Build finalizes the context creation and returns a new context.Context
// containing all provided fields and a cancel function with the given timeout.
//
// Internally, it uses context.WithValue for each key and caches all key names
// under an internal key so they can later be retrieved via GetCtxAll.
func (b CtxBuilder) Build() (context.Context, context.CancelFunc) {
	base := emptyContext
	keys := make([]string, 0, len(b.fields))

	for k, v := range b.fields {
		ctxKey := getCtxKey(k)
		base = context.WithValue(base, ctxKey, v)
		keys = append(keys, k)
	}

	base = context.WithValue(base, internalKeysKey, keys)
	return context.WithTimeout(base, b.timeout)
}

// GetCtx retrieves the string value for the given key from the context.
// Returns an empty string if the context is nil or key is not found.
func GetCtx(ctx context.Context, keyName ...string) string {
	if ctx == nil {
		return ""
	}

	if len(keyName) == 0 || keyName[0] == "" {
		return ""
	}

	key := getCtxKey(keyName[0])
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}

// GetCtxAll returns a map[string]string with all fields previously injected
// into the context using the CtxBuilder. Returns nil if context is nil.
func GetCtxAll(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}

	result := make(map[string]string)

	rawKeys := ctx.Value(internalKeysKey)
	if keyList, ok := rawKeys.([]string); ok {
		for _, k := range keyList {
			val := ctx.Value(getCtxKey(k))
			if strVal, ok := val.(string); ok {
				result[k] = strVal
			}
		}
	}

	return result
}
