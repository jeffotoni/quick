package recover

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/jeffotoni/quick"
)

// Config represents the configuration for the Recover middleware.
type Config struct {
	// The Quick instance to apply the middleware to.
	App *quick.Quick

	// Next is an optional function that, if returns true, skips the middleware.
	// Useful to conditionally bypass the recover logic for certain requests.
	Next func(c *quick.Ctx) bool

	// EnableStacktrace defaults to true. If false, the stacktrace will not be printed.
	// If App is not provided, the function panics, as it is required.
	EnableStacktrace bool
}

// New creates a new Recover middleware instance.
func New(config ...Config) func(next quick.Handler) quick.Handler {
	return func(next quick.Handler) quick.Handler {
		cfg := defaultConfig(config...)

		return quick.HandlerFunc(func(c *quick.Ctx) error {

			// Skip middleware if Next returns true
			if cfg.Next != nil && cfg.Next(c) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Fprint(os.Stderr, "Skipped recover middleware")
					}
				}()
				return next.ServeQuick(c)
			}

			// If stacktrace is disabled, send simple http error
			if !cfg.EnableStacktrace && cfg.App != nil {
				fmt.Fprintln(os.Stderr, "Warning: Stacktrace is disabled. Enable it for better debugging.")

				return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			}

			// If stacktrace is enabled, handle panics and return a 500 error
			if cfg.App != nil && cfg.EnableStacktrace {
				defer func() {
					if r := recover(); r != nil {
						// convert panic value to error
						err, ok := r.(error)
						if !ok {
							err = fmt.Errorf("%v", r)
						}
						fmt.Fprintf(os.Stdout, "Recovered from Quick.App: %v\n%s\n", err, debug.Stack())

						// Send a response to the client
						c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
					}
				}()
				return next.ServeQuick(c)
			}

			// serve the next handler
			return next.ServeQuick(c)
		})
	}
}

// defaultConfig returns the default configuration for the Recover middleware.
// If no configuration is provided, it returns a default configuration.
func defaultConfig(config ...Config) Config {
	if len(config) == 0 {
		// Default configuration
		return Config{
			App: nil,
			Next: func(c *quick.Ctx) bool {
				return false
			},
			EnableStacktrace: true,
		}
	}

	c := config[0]

	// halt the server
	if c.App == nil {
		panic("Config.App (Quick instance) is required")
	}

	// set default Next function
	if c.Next == nil {
		c.Next = func(c *quick.Ctx) bool {
			return false
		}
	}

	return c
}
