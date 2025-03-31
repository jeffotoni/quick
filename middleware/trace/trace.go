// Package trace provides middleware for the Quick framework that injects
// values into the Go context (`context.Context`) during request handling.
//
// This is useful for storing trace metadata like `X-Trace-ID`, `X-User-ID`,
// and environment tags (`env`, etc.) which can be accessed later during
// request lifecycle (e.g., for logging, telemetry, or business logic).
//
// Values are stored using private keys to avoid collision and can be safely
// retrieved using `GetCtx` or `GetCtxMap`.
//
// Usage Example:
//
//		app.Use(trace.New(trace.Config{
//			Timeout: 10 * time.Second,
//			Fields: map[string]func(c *quick.Ctx) string{
//				"X-Trace-ID": func(c *quick.Ctx) string {
//					traceID := c.Get("X-Trace-ID")
//					if traceID == "" {
//						traceID = rand.TraceID()
//					}
//					c.Set("X-Trace-ID", traceID)
//					return traceID
//				},
//				"X-User-ID": func(c *quick.Ctx) string {
//					userID := c.Get("X-User-ID")
//					if userID == "" {
//						userID = rand.AlgoDefault(1000, 9999)
//					}
//					c.Set("X-User-ID", userID)
//					return userID
//				},
//				"env": func(c *quick.Ctx) string {
//					return "dev"
//				},
//			},
//		}))
//
//	 After setup, you can access values inside any handler:
//
//		traceID := trace.GetCtx(c.Context(), "X-Trace-ID")
//		values := trace.GetCtxMap(c.Context())
package trace

import (
	"context"
	"sync"
	"time"

	"github.com/jeffotoni/quick"
)

type contextKey struct{ name string }

const internalCtxKeysKey = "__trace_keys__"

var keyCache sync.Map

// Config holds middleware settings for injecting context.
type Config struct {
	// Fields to inject into context (key + function to extract from *quick.Ctx)
	Fields map[string]func(*quick.Ctx) string

	// Timeout defines context timeout duration (default: 30s)
	Timeout time.Duration

	// Next skips the middleware if it returns true
	Next func(*quick.Ctx) bool
}

// getCtxKey returns a pointer-safe key for context usage.
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

// New returns the middleware to inject context values.
func New(cfg Config) func(next quick.Handler) quick.Handler {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	if cfg.Next == nil {
		cfg.Next = func(c *quick.Ctx) bool { return false }
	}

	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			if cfg.Next(c) {
				return next.ServeQuick(c)
			}

			base := context.Background()
			var keysUsed []string

			for key, fn := range cfg.Fields {
				if fn == nil {
					continue
				}
				val := fn(c)
				if val != "" {
					base = context.WithValue(base, getCtxKey(key), val)
					keysUsed = append(keysUsed, key)
				}
			}

			// Inject keys used for later GetCtxMap
			base = context.WithValue(base, getCtxKey(internalCtxKeysKey), keysUsed)

			ctx, cancel := context.WithTimeout(base, cfg.Timeout)
			defer cancel()

			// Store in quick context
			c.SetCtx(ctx)

			return next.ServeQuick(c)
		})
	}
}

// GetCtx retrieves a value from context by key.
func GetCtx(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	val := ctx.Value(getCtxKey(key))
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// GetCtxMap retrieves all key-value pairs injected by the middleware.
func GetCtxMap(ctx context.Context) map[string]string {
	result := make(map[string]string)
	if ctx == nil {
		return result
	}
	raw := ctx.Value(getCtxKey(internalCtxKeysKey))
	names, ok := raw.([]string)
	if !ok {
		return result
	}
	for _, name := range names {
		val := ctx.Value(getCtxKey(name))
		if str, ok := val.(string); ok {
			result[name] = str
		}
	}
	return result
}
