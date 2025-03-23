package healthcheck

import (
	"net/http"

	"github.com/jeffotoni/quick"
)

// Options defines the configuration for the healthcheck middleware.
// It allows customization of the endpoint, health probe logic, and conditional skipping.
type Options struct {
	// Next is an optional function that, if returns true, skips the middleware.
	// Useful to conditionally bypass the healthcheck logic for certain requests.
	Next func(c *quick.Ctx) bool

	// Endpoint specifies the route path that will be registered for the healthcheck.
	// Default: "/healthcheck"
	Endpoint string

	// Probe is a function executed during healthcheck requests.
	// It should return true if the application is healthy, false otherwise.
	// Default: always returns true.
	Probe func(c *quick.Ctx) bool

	// App is the instance of the Quick application.
	// Required to register the healthcheck endpoint once during setup.
	App *quick.Quick
}

// New initializes the healthcheck middleware and registers the health endpoint.
//
// Example usage:
//
//	app := quick.New()
//	app.Use(healthcheck.New(healthcheck.Options{
//		App:      app,
//		Endpoint: "/health",
//		Probe: func(c *quick.Ctx) bool {
//			// Custom health logic here
//			return true
//		},
//	}))
func New(opt ...Options) func(next quick.Handler) quick.Handler {
	option := defaultOptions(opt...)

	// Register the healthcheck route once, during app setup
	option.App.Get(option.Endpoint, func(c *quick.Ctx) error {
		// Skip route logic if Next returns true
		if option.Next != nil && option.Next(c) {
			return c.Status(http.StatusNotFound).SendString("Not Found")
		}

		// Only allow GET requests
		if c.Method() != quick.MethodGet {
			return c.Status(http.StatusMethodNotAllowed).SendString("Method Not Allowed")
		}

		// Execute health probe
		if option.Probe(c) {
			return c.Status(http.StatusOK).SendString("OK")
		}

		return c.Status(http.StatusServiceUnavailable).SendString("Service Unavailable")
	})

	// This middleware does not alter the request flow; it simply forwards it
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			return next.ServeQuick(c)
		})
	}
}

// defaultOptions applies sane defaults for the healthcheck middleware.
// If Endpoint or Probe are not provided, they are initialized to defaults.
// If App is not provided, the function panics, as it is required.
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
		panic("healthcheck.New: Options.App (Quick instance) is required to register the endpoint")
	}

	return cfg
}
