// Package pprof provides a Quick middleware that integrates Go's built-in
// net/http/pprof profiler for runtime analysis, performance debugging,
// and memory/cpu profiling.
//
// It exposes profiling endpoints such as /debug/pprof/, /heap, /goroutine, /profile, etc.
//
// Security:
//
// This middleware is intended for use in development environments.
// In production, it is recommended to disable or restrict access
// due to potential performance impact and exposure of internal details.
package pprof

import (
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/jeffotoni/quick"
)

// Config defines options for the pprof middleware.
//
// Prefix sets the base route for the pprof endpoints (default: "/debug/pprof").
// Next allows conditional execution, letting you disable the middleware
// for specific requests or environments (e.g., production).
type Config struct {
	// Prefix defines the base route for pprof endpoints.
	// Default is "/debug/pprof"
	Prefix string

	// Next is a function that allows skipping this middleware conditionally.
	// If it returns true, the middleware is bypassed.
	Next func(c *quick.Ctx) bool
}

// defaultConfig returns the default configuration for the middleware.
var defaultConfig = Config{
	Prefix: "/debug/pprof",
	Next:   nil,
}

// New returns a middleware handler for the Quick framework that serves
// pprof profiling data under a configurable prefix.
//
// The middleware automatically handles paths like /heap, /profile, /goroutine, etc.,
// using Go's net/http/pprof package.
//
// Optionally, a Config can be passed to set a custom prefix or control execution
// (e.g., disable in production).
//
// Example:
//
//	q.Use(pprof.New(pprof.Config{
//		Prefix: "/debug/pprof",
//		Next: func(c *quick.Ctx) bool {
//			return os.Getenv("APP_ENV") == "production"
//		},
//	}))
func New(config ...Config) func(next quick.Handler) quick.Handler {
	cfg := defaultConfig
	if len(config) > 0 {
		cfg = config[0]
		if cfg.Prefix == "" {
			cfg.Prefix = defaultConfig.Prefix
		}
	}

	// Middleware handler that intercepts requests matching the configured pprof prefix.
	// It serves profiling endpoints like /cmdline, /profile, /symbol, etc.
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Skip middleware if Next function is defined and returns true
			if cfg.Next != nil && cfg.Next(c) {
				return next.ServeQuick(c)
			}

			path := c.Path()

			// Check if the request path matches the pprof prefix
			if path == cfg.Prefix || strings.HasPrefix(path, cfg.Prefix+"/") {
				// Trim the prefix to get the specific subpath (e.g., /cmdline, /heap)
				subpath := strings.TrimPrefix(path, cfg.Prefix)

				// If subpath is empty or root, serve the pprof index
				if subpath == "" || subpath == "/" {
					http.HandlerFunc(pprof.Index).ServeHTTP(c.Response, c.Request)
					return nil
				}

				// Match each supported pprof route and serve accordingly
				switch subpath {
				case "/cmdline":
					http.HandlerFunc(pprof.Cmdline).ServeHTTP(c.Response, c.Request)
				case "/profile":
					http.HandlerFunc(pprof.Profile).ServeHTTP(c.Response, c.Request)
				case "/symbol":
					http.HandlerFunc(pprof.Symbol).ServeHTTP(c.Response, c.Request)
				case "/trace":
					http.HandlerFunc(pprof.Trace).ServeHTTP(c.Response, c.Request)
				case "/allocs":
					http.HandlerFunc(pprof.Handler("allocs").ServeHTTP).ServeHTTP(c.Response, c.Request)
				case "/block":
					http.HandlerFunc(pprof.Handler("block").ServeHTTP).ServeHTTP(c.Response, c.Request)
				case "/goroutine":
					http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP).ServeHTTP(c.Response, c.Request)
				case "/heap":
					http.HandlerFunc(pprof.Handler("heap").ServeHTTP).ServeHTTP(c.Response, c.Request)
				case "/mutex":
					http.HandlerFunc(pprof.Handler("mutex").ServeHTTP).ServeHTTP(c.Response, c.Request)
				case "/threadcreate":
					http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP).ServeHTTP(c.Response, c.Request)
				default:
					// Redirect unknown subpaths to the index page
					return c.Redirect(cfg.Prefix + "/")
				}

				return nil
			}

			// If path doesn't match, continue to the next middleware or handler
			return next.ServeQuick(c)
		})
	}
}
