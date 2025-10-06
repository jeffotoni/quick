// Package recover provides a middleware for the Quick web framework
// that gracefully handles panics during request processing.
//
// When a panic occurs, the middleware intercepts it, optionally prints a stack trace,
// and sends a 500 Internal Server Error response to the client.
//
// You can customize the behavior using the Config struct:
//   - Enable or disable stack trace logging
//   - Provide a custom stack trace handler for advanced logging or reporting
//   - Conditionally skip the middleware using a Next function
//
// This middleware helps ensure that your application remains available and does not crash
// due to unexpected panics in route handlers or other middlewares.
package recover

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jeffotoni/quick"
)

// Config defines the configuration for the Recover middleware.
type Config struct {
	// Next is an optional function. If it returns true, the middleware is skipped.
	// This is useful to conditionally disable panic recovery for specific requests.
	Next func(c *quick.Ctx) bool

	// EnableStacktrace enables printing the stack trace to stderr when a panic occurs.
	// Defaults to true.
	EnableStacktrace bool

	// StackTraceHandler is an optional function that handles the recovered panic.
	// If set, it will be called instead of the default stack trace printer.
	// Useful for custom logging or error reporting systems.
	StackTraceHandler func(c *quick.Ctx, err interface{})
}

// New returns a Recover middleware that captures panics during request handling.
// It recovers from the panic, logs the error (using either the default or custom handler),
// and responds with a 500 Internal Server Error.
func New(cfgs ...Config) func(next quick.Handler) quick.Handler {
	// Apply default config if none is provided

	cfg := defaultConfig(cfgs...)

	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Defer a function to recover from any panics that occur during request handling
			defer func() {
				if r := recover(); r != nil {
					handlePanic(r, cfg, c)
				}
			}()

			// Skip middleware logic if Next returns true
			if cfg.Next != nil && cfg.Next(c) {
				return next.ServeQuick(c)
			}

			// Continue with the next handler in the chain
			return next.ServeQuick(c)
		})
	}
}

// defaultConfig applies sensible defaults for the Recover middleware.
// If no config is provided, it initializes default values.
func defaultConfig(config ...Config) Config {
	// No config passed, return default with stack trace enabled
	if len(config) == 0 {
		return Config{
			EnableStacktrace: true,
		}
	}

	c := config[0]

	// Set default behavior for Next if not defined
	if c.Next == nil {
		c.Next = func(c *quick.Ctx) bool {
			return false
		}
	}

	return c
}

// handlePanic processes a recovered panic.
// If StackTraceHandler is defined, it delegates to that handler.
// Otherwise, it logs the panic to stderr with or without a stack trace,
// and sends a 500 Internal Server Error response to the client.
func handlePanic(r interface{}, cfg Config, c *quick.Ctx) {
	// Convert recovered value to error type if necessary
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}

	// Use the custom stack trace handler if provided.
	if cfg.StackTraceHandler != nil {
		cfg.StackTraceHandler(c, r)
	} else if cfg.EnableStacktrace {
		// Print error with stack trace
		fmt.Fprintf(os.Stderr, "Recovered panic: %v\n%s\n", err, debug.Stack())
	} else {
		// Print error without stack trace
		fmt.Fprintln(os.Stderr, "Recovered panic: stacktrace disabled.")
	}

	if c != nil && c.Response != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to send error response: %v\n", err)
			}
		}()
	}

	// Respond with HTTP 500 Internal Server Error
	c.Status(quick.StatusInternalServerError).SendString("Internal Server Error")
}
