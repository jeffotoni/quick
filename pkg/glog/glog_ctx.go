package glog

import (
	"context"
	"sync"
	"time"
)

// contextKey is a private type to avoid collisions in context
type contextKey struct{ name string }

const defaultCtxKey = "TraceID"
const defaultCtxTimeout = 30 * time.Second

var keyCache sync.Map // map[string]*contextKey

func getCtxKey(name string) *contextKey {
	if name == "" {
		name = defaultCtxKey
	}
	if v, ok := keyCache.Load(name); ok {
		return v.(*contextKey)
	}
	k := &contextKey{name}
	keyCache.Store(name, k)
	return k
}

// CtxBuilder provides a fluent API to build a context with a trace ID.
type CtxBuilder struct {
	name    string
	key     string
	timeout time.Duration
}

// NewCtx creates a new fluent context builder.
func NewCtx() *CtxBuilder {
	return &CtxBuilder{
		name:    defaultCtxKey,
		timeout: defaultCtxTimeout,
	}
}

// Name sets a custom context key name.
func (b *CtxBuilder) Name(name string) *CtxBuilder {
	if name != "" {
		b.name = name
	}
	return b
}

// Key sets the value to store in the context.
func (b *CtxBuilder) Key(val string) *CtxBuilder {
	b.key = val
	return b
}

// Timeout sets a custom timeout duration for the context.
func (b *CtxBuilder) Timeout(d time.Duration) *CtxBuilder {
	if d > 0 {
		b.timeout = d
	}
	return b
}

// Build creates the context and returns it with a cancel function.
func (b *CtxBuilder) Build() (context.Context, context.CancelFunc) {
	ctxKey := getCtxKey(b.name)
	base := context.WithValue(context.Background(), ctxKey, b.key)
	ctx, cancel := context.WithTimeout(base, b.timeout)
	return ctx, cancel
}

// GetCtx retrieves the trace ID from context using the given key name (optional).
// Defaults to "TraceID" if no keyName is provided.
func GetCtx(ctx context.Context, keyName ...string) string {
	if ctx == nil {
		return ""
	}
	key := defaultCtxKey
	if len(keyName) > 0 && keyName[0] != "" {
		key = keyName[0]
	}
	ctxKey := getCtxKey(key)

	if val, ok := ctx.Value(ctxKey).(string); ok {
		return val
	}
	return ""
}
