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
	"os"

	"github.com/jeffotoni/quick"
)

// Options defines the configuration for the profiling middleware.
type Options struct {
	// Next is an optional function that, if returns true, skips the middleware.
	// Useful to conditionally bypass the middleware.
	Next func(c *quick.Ctx) bool

	// App is the instance of the Quick application.
	// Required to register the profiling endpoint once during setup.
	App *quick.Quick
}

func New(opt ...Options) func(next quick.Handler) quick.Handler {
	option := defaultOptions(opt...)

	// Define profiling routes
	profilingRoutes := map[string]http.Handler{
		"/debug/pprof":        http.HandlerFunc(pprof.Index),
		"/debug/cmdline":      http.HandlerFunc(pprof.Cmdline),
		"/debug/profile":      http.HandlerFunc(pprof.Profile),
		"/debug/symbol":       http.HandlerFunc(pprof.Symbol),
		"/debug/pprof/trace":  http.HandlerFunc(pprof.Trace),
		"/debug/goroutine":    pprof.Handler("goroutine"),
		"/debug/heap":         pprof.Handler("heap"),
		"/debug/threadcreate": pprof.Handler("threadcreate"),
		"/debug/mutex":        pprof.Handler("mutex"),
		"/debug/allocs":       pprof.Handler("allocs"),
		"/debug/block":        pprof.Handler("block"),
	}

	// Register routes in Quick
	for route, handler := range profilingRoutes {
		option.App.Get(route, func(c *quick.Ctx) error {
			// Skip route logic if Next returns true
			if option.Next != nil && option.Next(c) {
				return c.Status(http.StatusNotFound).SendString("Not Found")
			}

			// Only allow GET requests
			if c.Method() != quick.MethodGet {
				return c.Status(http.StatusMethodNotAllowed).SendString("Method Not Allowed")
			}

			// Serve the profiling endpoint
			handler.ServeHTTP(c.Response, c.Request)
			return nil
		})
	}

	// Middleware just forwards the request
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			return next.ServeQuick(c)
		})
	}
}

// defaultOptions applies sane defaults for the profiling middleware.
// If App is not provided, the function panics, as it is required.
// Profiling is only enabled in development mode
func defaultOptions(opt ...Options) Options {
	// Check if APP_ENV is set to "development".
	// Profiling is only enabled in development mode
	env := os.Getenv("APP_ENV")
	if env != "development" {
		panic("pprof.New: Environment variable APP_ENV must be set to 'development'")
	}

	// Initialize with default values
	if len(opt) == 0 {
		return Options{
			Next: func(c *quick.Ctx) bool { return false },
		}
	}

	// Check if App is provided
	if opt[0].App == nil {
		panic("pprof.New: Options.App (Quick instance) is required to register the profiling route")
	}
	return opt[0]
}
