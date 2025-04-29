package cache

import (
	"fmt"
	"time"

	"github.com/jeffotoni/quick"
)

// ExampleNew_defaultBehavior demonstrates the default behavior of the cache middleware.
// This function is named ExampleNew_defaultBehavior()
// it with the Examples type.
func ExampleNew_defaultBehavior() {
	q := quick.New()

	// Use the cache middleware with default configuration
	q.Use(New())

	// Define a route that returns the current time
	q.Get("/time", func(c *quick.Ctx) error {
		return c.String(fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC3339)))
	})

	// Make a request to the route
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time",
	})

	// First request should miss the cache
	fmt.Println("First request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make a second request to the same route
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time",
	})

	// Second request should hit the cache
	fmt.Println("Second request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Output:
	// First request cache status: MISS
	// Second request cache status: HIT
}

// ExampleNew_customKeyGenerator demonstrates how to use a custom key generator.
// This function is named ExampleNew_customKeyGenerator()
// it with the Examples type.
func ExampleNew_customKeyGenerator() {
	q := quick.New()

	// Use the cache middleware with a custom key generator
	q.Use(New(Config{
		KeyGenerator: func(c *quick.Ctx) string {
			return c.Path() + "?lang=" + c.Query["lang"]
		},
	}))

	// Define a route that returns a greeting in the requested language
	q.Get("/greeting", func(c *quick.Ctx) error {
		lang := c.Query["lang"]
		greeting := "Hello"
		if lang == "es" {
			greeting = "Hola"
		} else if lang == "fr" {
			greeting = "Bonjour"
		}
		return c.String(fmt.Sprintf("%s, World! (Generated at %s)", greeting, time.Now().Format(time.RFC3339)))
	})

	// Make a request with lang=en
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/greeting?lang=en",
	})
	fmt.Println("English request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make a request with lang=es
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/greeting?lang=es",
	})
	fmt.Println("Spanish request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make another request with lang=en (should hit cache)
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/greeting?lang=en",
	})
	fmt.Println("Second English request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Output:
	// English request cache status: MISS
	// Spanish request cache status: MISS
	// Second English request cache status: HIT
}

// ExampleNew_cacheInvalidation demonstrates how to use cache invalidation.
// This function is named ExampleNew_cacheInvalidation()
// it with the Examples type.
func ExampleNew_cacheInvalidation() {
	q := quick.New()

	// Use the cache middleware with a cache invalidator
	q.Use(New(Config{
		CacheInvalidator: func(c *quick.Ctx) bool {
			return c.Query["clear"] == "1"
		},
	}))

	// Define a route that returns the current time
	q.Get("/time", func(c *quick.Ctx) error {
		return c.String(fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC3339)))
	})

	// Make a request to the route
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time",
	})
	fmt.Println("First request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make a second request (should hit cache)
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time",
	})
	fmt.Println("Second request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make a request with clear=1 (should invalidate cache)
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time?clear=1",
	})
	fmt.Println("Invalidation request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Make another request (should miss cache)
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/time",
	})
	fmt.Println("Post-invalidation request cache status:", resp.Response().Header.Get("X-Cache-Status"))

	// Output:
	// First request cache status: MISS
	// Second request cache status: HIT
	// Invalidation request cache status: INVALIDATED
	// Post-invalidation request cache status: MISS
}

// ExampleNew_customExpiration demonstrates how to use custom expiration times.
// This function is named ExampleNew_customExpiration()
// it with the Examples type.
func ExampleNew_customExpiration() {
	q := quick.New()

	// Use the cache middleware with a custom expiration generator
	q.Use(New(Config{
		ExpirationGenerator: func(c *quick.Ctx, cfg *Config) time.Duration {
			if c.Path() == "/short" {
				return 1 * time.Second
			}
			return cfg.Expiration // Use default for other paths
		},
	}))

	// Define routes with different cache durations
	q.Get("/short", func(c *quick.Ctx) error {
		return c.String(fmt.Sprintf("Short cache: %s", time.Now().Format(time.RFC3339)))
	})

	q.Get("/default", func(c *quick.Ctx) error {
		return c.String(fmt.Sprintf("Default cache: %s", time.Now().Format(time.RFC3339)))
	})

	// Make requests to both routes
	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/short",
	})
	fmt.Println("Short cache first request:", resp.Response().Header.Get("X-Cache-Status"))

	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/default",
	})
	fmt.Println("Default cache first request:", resp.Response().Header.Get("X-Cache-Status"))

	// Make second requests (both should hit cache)
	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/short",
	})
	fmt.Println("Short cache second request:", resp.Response().Header.Get("X-Cache-Status"))

	resp, _ = q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/default",
	})
	fmt.Println("Default cache second request:", resp.Response().Header.Get("X-Cache-Status"))

	// Output:
	// Short cache first request: MISS
	// Default cache first request: MISS
	// Short cache second request: HIT
	// Default cache second request: HIT
}
