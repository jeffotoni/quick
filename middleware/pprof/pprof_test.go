package pprof

import (
	"os"
	"testing"

	"github.com/jeffotoni/quick"
)

func TestWithDefaultConfig(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	q := quick.New()
	q.Use(New(Options{
		App: q,
	}))

	// Define a test route
	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/debug/pprof",
	})

	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error(err)
	}
	os.Unsetenv("APP_ENV")
}

func TestWithRoutePrefix(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	q := quick.New()
	q.Use(New(Options{
		App: q,
	}))

	// Define a test route
	routes := []string{
		"/debug/cmdline",
		"/debug/profile",
		"/debug/symbol",
		"/debug/pprof/trace",
		"/debug/goroutine",
		"/debug/heap",
		"/debug/threadcreate",
		"/debug/mutex",
		"/debug/allocs",
		"/debug/block",
	}

	for _, route := range routes {
		t.Run("Testing "+route, func(t *testing.T) {
			resp, err := q.Qtest(quick.QuickTestOptions{
				Method: quick.MethodGet,
				URI:    route,
			})

			if err != nil {
				t.Fatal(err)
			}
			if err := resp.AssertStatus(quick.StatusOK); err != nil {
				t.Errorf("Route %s: %v", route, err)
			}
		})
	}
	os.Unsetenv("APP_ENV")
}
