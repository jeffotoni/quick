// Package glog provides a lightweight, flexible, and structured logging library for Go.
// It supports multiple formats (text, slog-style, and JSON), dynamic fields, custom patterns,
// caller tracing, and context propagation for values like TraceID, X-User-ID, etc.
//
// It is part of the Quick Framework: https://github.com/jeffotoni/quick
//
// ## Features:
//
//   - Text, JSON, and slog-style logging
//   - Fluent API and legacy T-style (`InfoT`, `DebugT`) support
//   - Customizable output patterns with placeholders (${time}, ${msg}, etc)
//   - Dynamic field injection with ordered rendering
//   - Context helpers for setting and retrieving traceable values (e.g., TraceID)
//   - Full color support in terminal for levels and values
//   - Custom log levels: DEBUG, INFO, WARN, ERROR
//   - Global fields (`CustomFields`) and per-call contextual fields
//
// ## Example:
//
//	func main() {
//	    glog.Set(glog.Config{
//	        Format:     "text",
//	        Pattern:    "[${time}] ${level} ${msg} |",
//	        Level:      glog.DEBUG,
//	        TimeFormat: time.RFC3339,
//	    })
//
//	    glog.Info("App started").
//	        Str("version", "1.0.0").
//	        Str("env", "production").
//	        Send()
//
//	    ctx, cancel := glog.NewCtx().
//	        Set("TraceID", "abc-123").
//	        Set("X-User-ID", "user-789").
//	        Timeout(5 * time.Second).
//	        Build()
//	    defer cancel()
//
//	    traceID := glog.GetCtx(ctx)
//	    userID := glog.GetCtx(ctx, "X-User-ID")
//
//	    glog.Debug("Request received").
//	        Str("trace", traceID).
//	        Str("user", userID).
//	        Send()
//	}
//
// Output (text):
//
//	[2025-03-30T15:20:00Z] INFO App started | version 1.0.0 env production
//	[2025-03-30T15:20:00Z] DEBUG Request received | trace abc-123 user user-789
package glog

import (
	"context"
	"sync"
	"time"
)

const internalCtxKeysKey = "__glog_keys__"

// contextKey is a private type to avoid collisions in context
type contextKey struct{ name string }

const defaultCtxTimeout = 30 * time.Second

var keyCache sync.Map // map[string]*contextKey

// CtxBuilder provides a fluent API to build a context with one or multiple keys.
type CtxBuilder struct {
	values   map[string]string
	keysUsed []string
	timeout  time.Duration
}

// NewCtx creates a new fluent context builder.
func NewCtx() *CtxBuilder {
	return &CtxBuilder{
		values:  make(map[string]string),
		timeout: defaultCtxTimeout,
	}
}

// Set injects a value into the context under a given name.
func (b *CtxBuilder) Set(name, value string) *CtxBuilder {
	if name != "" && value != "" {
		if _, exists := b.values[name]; !exists {
			b.keysUsed = append(b.keysUsed, name)
		}
		b.values[name] = value
	}
	return b
}

// Timeout sets a custom timeout duration for the context.
func (b *CtxBuilder) Timeout(d time.Duration) *CtxBuilder {
	if d > 0 {
		b.timeout = d
	}
	return b
}

// getCtxKey returns a unique pointer-based context key for the given name.
// It ensures consistent key usage by caching pointers for each key name,
// preventing accidental key collisions.
//
// If the name is empty, it returns a pointer to a default key with the name "TraceID".
//
// Internally uses sync.Map to avoid recreating keys and to enable identity-based lookups
// (ensuring different keys don't conflict even if their string values are the same).
func getCtxKey(name string) *contextKey {
	if name == "" {
		return &contextKey{"TraceID"}
	}
	if v, ok := keyCache.Load(name); ok {
		return v.(*contextKey)
	}
	k := &contextKey{name}
	keyCache.Store(name, k)
	return k
}

func (b *CtxBuilder) Build() (context.Context, context.CancelFunc) {
	base := context.Background()
	for name, val := range b.values {
		base = context.WithValue(base, getCtxKey(name), val)
	}
	// Store the used keys slice inside the context
	base = context.WithValue(base, getCtxKey(internalCtxKeysKey), b.keysUsed)
	return context.WithTimeout(base, b.timeout)
}

// GetCtx retrieves a single value from context using a given key.
// Defaults to "TraceID" if no keyName is provided.
func GetCtx(ctx context.Context, keyName ...string) string {
	if ctx == nil {
		return ""
	}
	key := "TraceID"
	if len(keyName) > 0 && keyName[0] != "" {
		key = keyName[0]
	}
	ctxKey := getCtxKey(key)
	if val, ok := ctx.Value(ctxKey).(string); ok {
		return val
	}
	return ""
}

// GetCtxMap retrieves all known string keys injected using Set().
func GetCtxMap(ctx context.Context) map[string]string {
	result := make(map[string]string)
	keysRaw := ctx.Value(getCtxKey(internalCtxKeysKey))
	keys, ok := keysRaw.([]string)
	if !ok {
		return result
	}
	for _, name := range keys {
		key := getCtxKey(name)
		if val, ok := ctx.Value(key).(string); ok {
			result[name] = val
		}
	}
	return result
}
