package pprof

import (
	"os"
	"strings"
	"testing"

	"github.com/jeffotoni/quick"
)

// TestPprofIndex verifies that the main pprof index route (/debug/pprof)
// returns a 200 OK response when the pprof middleware is enabled.
// This ensures that the profiler is properly registered in the application.
//
// To run:
//
//	go test -v -run ^TestPprofIndex$
func TestPprofIndex(t *testing.T) {
	// Set environment to development to enable the pprof middleware
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	// Initialize the Quick application
	app := quick.New()

	// Register the pprof middleware
	app.Use(New())

	// Declare the route pattern so pprof middleware can intercept requests
	app.Get("/debug/pprof*", func(c *quick.Ctx) error { return nil })

	// Perform a test request to the pprof index route
	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/debug/pprof",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the response status is 200 OK
	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error("Expected 200 OK for /debug/pprof, got:", err)
	}
}

// TestPprofRoutes verifies that all standard pprof subroutes return 200 OK.
// This ensures each profiling endpoint is correctly exposed by the middleware.
//
// To run:
//
//	go test -v -run ^TestPprofRoutes$
func TestPprofRoutes(t *testing.T) {
	// Set environment to development to enable the pprof middleware
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	// Initialize the Quick application
	app := quick.New()

	// Register the pprof middleware
	app.Use(New())

	// Declare a wildcard route for the pprof endpoints
	app.Get("/debug/pprof*", func(c *quick.Ctx) error { return nil })

	// List of pprof subroutes to test
	routes := []string{
		"/debug/pprof/cmdline",
		"/debug/pprof/profile?seconds=1",
		"/debug/pprof/symbol",
		"/debug/pprof/trace",
		"/debug/pprof/goroutine",
		"/debug/pprof/heap",
		"/debug/pprof/threadcreate",
		"/debug/pprof/mutex",
		"/debug/pprof/allocs",
		"/debug/pprof/block",
	}

	// Run each route and assert status 200 OK
	for _, route := range routes {
		t.Run("Testing "+route, func(t *testing.T) {
			resp, err := app.Qtest(quick.QuickTestOptions{
				Method: quick.MethodGet,
				URI:    route,
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := resp.AssertStatus(quick.StatusOK); err != nil {
				t.Errorf("Expected 200 OK for %s, got: %v", route, err)
			}
		})
	}
}

// TestPprofInvalidPathRedirect ensures that accessing an invalid
// pprof subroute triggers a redirect back to the index page (/debug/pprof/).
//
// To run:
//
//	go test -v -run ^TestPprofInvalidPathRedirect$
func TestPprofInvalidPathRedirect(t *testing.T) {
	// Set environment to development to enable the pprof middleware
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	// Initialize the Quick application
	app := quick.New()

	// Register the pprof middleware
	app.Use(New())

	// Declare a wildcard route for the pprof endpoints
	app.Get("/debug/pprof*", func(c *quick.Ctx) error { return nil })

	// Send request to an unknown pprof subpath (not handled directly by pprof)
	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/debug/pprof/unknown",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the response is a redirect (HTTP 302 Found)
	if err := resp.AssertStatus(quick.StatusFound); err != nil {
		t.Error("Expected redirect (302) for unknown pprof path, got:", err)
	}
}

// TestPprofWithCustomPrefixAndNext tests two features of the middleware:
// 1. Custom prefix configuration (e.g., /custom-pprof).
// 2. Conditional skipping via the Next function.
// It checks that profiling works under the custom prefix and that paths matching
// a specific condition are correctly excluded from profiling.
//
// To run:
//
//	go test -v -run ^TestPprofWithCustomPrefixAndNext$
func TestPprofWithCustomPrefixAndNext(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	app := quick.New()

	// Middleware with custom prefix and Next logic
	app.Use(New(Config{
		Prefix: "/custom-pprof",
		Next: func(c *quick.Ctx) bool {
			// Skip profiling if path contains "skip"
			return strings.Contains(c.Path(), "skip")
		},
	}))

	// Declare wildcard route to enable middleware execution
	app.Get("/custom-pprof*", func(c *quick.Ctx) error {
		return c.Next()
	})

	// Should be handled by pprof
	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/custom-pprof",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error("Expected 200 OK for /custom-pprof, got:", err)
	}

	// Should be skipped due to Next function returning true
	resp, err = app.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/custom-pprof/skip-this",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusNotFound); err != nil {
		t.Error("Expected 404 Not Found for skipped path, got:", err)
	}
}
