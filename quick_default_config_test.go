package quick

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestQuick_Listen this example tests whether the server can be started
// correctly and prevents multiple instances on the same port.
// The will test TestQuick_Listen(t *testing.T)
//
//	$ go test -v -run ^TestQuick_Listen
func TestQuick_Listen(t *testing.T) {

	// Other tests omitted for brevity...
	t.Run("Error trying to run server on the same port", func(t *testing.T) {
		q1 := New()
		server1, shutdown1, err1 := q1.ListenWithShutdown(":0")
		if err1 != nil {
			t.Fatalf("Unexpected error starting first server: %v", err1)
		}
		defer shutdown1()

		q2 := New()
		_, shutdown2, err2 := q2.ListenWithShutdown(server1.Addr)
		if err2 == nil {
			shutdown2()
			t.Errorf("Expected error running server on the same port (%s), but no error occurred", server1.Addr)
		} else {
			fmt.Println("Error when trying to run second server on the same port detected correctly.")
		}
	})
}

// TestQuick_ServeHTTP  Ensure that a recorded route responds correctly.
// The will test TestQuick_ServeHTTP(t *testing.T)
//
//	$ go test -v -run ^TestQuick_ServeHTTP
func TestQuick_ServeHTTP(t *testing.T) {
	q := New()

	// Register a test route
	q.Get("/ping", func(c *Ctx) error {
		return c.String("pong")
	})

	// Create a test server
	ts := httptest.NewServer(q)
	defer ts.Close()

	t.Run("Registered route responds correctly", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/ping")
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, but got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "pong" {
			t.Errorf("Expected 'pong' response, but got '%s'", body)
		}
	})

	t.Run("Unregistered route returns 404", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/does not exist")
		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, but got %d", resp.StatusCode)
		}
	})
}

// TestDefaultConfig verify that the default configuration of
// Quick is being initialized correctly with the expected values.
// The will test TestDefaultConfig(t *testing.T)
//
//	$ go test -v -run ^TestDefaultConfig
func TestDefaultConfig(t *testing.T) {
	expectedConfig := Config{
		BodyLimit:         2 * 1024 * 1024,
		MaxBodySize:       2 * 1024 * 1024,
		MaxHeaderBytes:    1 * 1024 * 1024,
		RouteCapacity:     1000,
		MoreRequests:      290,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
	}

	if defaultConfig != expectedConfig {
		t.Errorf("esperado %+v, mas obteve %+v", expectedConfig, defaultConfig)
	}
}

// TestQuickInitializationWithCustomConfig  Ensure that Quick can be booted with custom settings.
// The will test TestQuickInitializationWithCustomConfig(t *testing.T)
//
//	$ go test -v -run ^TestQuickInitializationWithCustomConfig
func TestQuickInitializationWithCustomConfig(t *testing.T) {
	customConfig := Config{
		BodyLimit:         4 * 1024 * 1024,
		MaxBodySize:       4 * 1024 * 1024,
		MaxHeaderBytes:    2 * 1024 * 1024,
		RouteCapacity:     500,
		MoreRequests:      500,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       2 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}

	q := New(customConfig)

	if q.config != customConfig {
		t.Errorf("esperado %+v, mas obteve %+v", customConfig, q.config)
	}
}

// TestQuickInitializationWithCustomConfig ensure that default values are applied
// correctly when creating an instance without explicit configuration.
// The will test func TestQuickInitializationDefaults(t *testing.T)
//
//	$ go test -v -run ^func TestQuickInitializationDefaults(t *testing.T)
func TestQuickInitializationDefaults(t *testing.T) {
	q := New()

	if q.config.BodyLimit != defaultConfig.BodyLimit {
		t.Errorf("BodyLimit incorreto: esperado %d, obteve %d", defaultConfig.BodyLimit, q.config.BodyLimit)
	}
	if q.config.MaxBodySize != defaultConfig.MaxBodySize {
		t.Errorf("MaxBodySize incorreto: esperado %d, obteve %d", defaultConfig.MaxBodySize, q.config.MaxBodySize)
	}
	if q.config.MoreRequests != defaultConfig.MoreRequests {
		t.Errorf("MoreRequests incorreto: esperado %d, obteve %d", defaultConfig.MoreRequests, q.config.MoreRequests)
	}
}

// TestQuickInitializationWithCustomConfig ensure that default values are applied
// correctly when creating an instance without explicit configuration.
// The will test func TestQuickInitializationDefaults(t *testing.T)
//
//	$ go test -v -run ^func TestQuickInitializationDefaults(t *testing.T)
func TestQuickInitializationWithZeroValues(t *testing.T) {
	zeroConfig := Config{}
	q := New(zeroConfig)

	if q.config.RouteCapacity != 1000 {
		t.Errorf("RouteCapacity incorreto: esperado 1000, obteve %d", q.config.RouteCapacity)
	}
}

// TestQuick_GetRoute ensure that routes are correctly recorded and retrieved
// correctly when creating an instance without explicit configuration.
// The will test func TestQuick_GetRoute(t *testing.T)
//
//	$ go test -v -run ^func TestQuick_GetRoute(t *testing.T)
func TestQuick_GetRoute(t *testing.T) {
	q := New()

	// Check if the route list is empty initially
	if len(q.GetRoute()) != 0 {
		t.Errorf("Expected 0 routes, but got %d", len(q.GetRoute()))
	}

	// Add a test route
	q.Get("/ping", func(c *Ctx) error {
		return c.String("pong")
	})

	// Check if the route was registered correctly
	routes := q.GetRoute()
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, but got %d", len(routes))
	}

	// Check if the route details are correct
	expectedPath := "/ping"
	if routes[0].Path != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, routes[0].Path)
	}

	expectedMethod := "GET"
	if routes[0].Method != expectedMethod {
		t.Errorf("Expected method '%s', but got '%s'", expectedMethod, routes[0].Method)
	}
}
