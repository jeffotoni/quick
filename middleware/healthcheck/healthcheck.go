package healthcheck

import (
	"github.com/jeffotoni/quick"
)

// Options defines the configuration for the Helmet middleware.
// All fields map to specific HTTP headers that enhance security.
// These options can override the default behavior provided by the middleware.
type Options struct {
	// Next defines a function to skip the middleware
	Next func(c *quick.Ctx) bool

	// Endpoint defines the path to the healthcheck endpoint
	Endpoint string

	// Probe defines the function to check the health of the application
	Probe func(*quick.Ctx) bool

	// App is the Quick instance
	App *quick.Quick
}

func New(opt ...Options) func(next quick.Handler) quick.Handler {
	// Apply default options
	option := defaultOptions(opt...)

	option.App.Get(option.Endpoint, func(c *quick.Ctx) error {
		if option.Next != nil && option.Next(c) {
			return c.Status(quick.StatusNotFound).String("Not Found")
		}

		if c.Method() != quick.MethodGet {
			return c.Status(quick.StatusMethodNotAllowed).String("Method Not Allowed")
		}

		if option.Probe(c) {
			return c.Status(quick.StatusOK).String("OK")
		}

		return c.Status(quick.StatusServiceUnavailable).String("Service Unavailable")
	})

	// The middleware itself can continue to pass through
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// The middleware does not interfere with the routes,
			// it just passes them on to the next one
			return next.ServeQuick(c)
		})
	}

}

// defaultOptions returns a set of default values for the healthcheck middleware.
func defaultOptions(opt ...Options) Options {
	cfg := Options{
		Endpoint: "/healthcheck",
		Probe: func(c *quick.Ctx) bool {
			return true
		},
	}

	if len(opt) > 0 {
		cfg = opt[0]
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "/healthcheck"
	}

	if cfg.Probe == nil {
		cfg.Probe = func(c *quick.Ctx) bool {
			return true
		}
	}

	if cfg.App == nil {
		panic("healthcheck.New: App instance must be provided in Options.App")
	}

	return cfg
}
