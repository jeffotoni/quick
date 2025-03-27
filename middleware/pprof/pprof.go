// Package profiling provides a middleware for profiling in Quick to help with
// debugging and optimization.
//
// It allows you to enable profiling in development mode only.
// In production, profiling is disabled, because it can introduce unwanted overhead
// and potentially degrade performance
package pprof

import (
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/jeffotoni/quick"
)

// Config defines the configuration for the pprof middleware.
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

// New returns a Quick middleware handler that exposes pprof endpoints
// at the configured prefix. It dynamically intercepts requests and serves
// profiling data without requiring manual route registration.
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
