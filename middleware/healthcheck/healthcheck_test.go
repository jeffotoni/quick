package healthcheck

import (
	"testing"

	"github.com/jeffotoni/quick"
)

// TestHealthcheck tests the healthcheck middleware with custom endpoint.
func TestHealthcheckWithCustomEndpoint(t *testing.T) {
	q := quick.New()
	q.Use(New(Options{
		Endpoint: "/v1/health",
		App:      q,
		Probe: func(c *quick.Ctx) bool {
			return true
		},
	}))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/health",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error(err)
	}
	if err := resp.AssertString("OK"); err != nil {
		t.Error(err)
	}
}

// TestHealthcheckEndpoint tests the healthcheck middleware with default endpoint.
func TestHealthcheckEndpoint(t *testing.T) {
	q := quick.New()
	q.Use(New(
		Options{
			App: q,
		},
	))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/healthcheck",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := resp.AssertStatus(quick.StatusOK); err != nil {
		t.Error(err)
	}

	if err := resp.AssertString("OK"); err != nil {
		t.Error(err)
	}
}

// TestHealthcheckProbeFalse tests when Probe returns false (service unavailable).
func TestHealthcheckProbeFalse(t *testing.T) {
	q := quick.New()
	q.Use(New(Options{
		App: q,
		Probe: func(c *quick.Ctx) bool {
			return false
		},
	}))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/healthcheck",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := resp.AssertStatus(quick.StatusServiceUnavailable); err != nil {
		t.Error(err)
	}
}

// TestHealthcheckMethodNotAllowed tests when method is not GET.
// func TestHealthcheckMethodNotAllowed(t *testing.T) {
// 	q := quick.New()
// 	q.Use(New(Options{
// 		App: q,
// 	}))

// 	resp, err := q.Qtest(quick.QuickTestOptions{
// 		Method: quick.MethodPost, // POST instead of GET
// 		URI:    "/healthcheck",
// 	})
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	if err := resp.AssertStatus(quick.StatusMethodNotAllowed); err != nil {
// 		t.Error(err)
// 	}
// }

// TestHealthcheckWithNextSkipping tests when Next() returns true (skips route logic).
func TestHealthcheckWithNextSkipping(t *testing.T) {
	q := quick.New()
	q.Use(New(Options{
		App: q,
		Next: func(c *quick.Ctx) bool {
			return true // Always skip
		},
	}))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/healthcheck",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := resp.AssertStatus(quick.StatusNotFound); err != nil {
		t.Error(err)
	}
}
