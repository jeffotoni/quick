package healthcheck

import (
	"testing"

	"github.com/jeffotoni/quick"
)

// TestHealthcheckWithCustomEndpoint verifies that the healthcheck middleware
//
// responds with status 200 and body "OK" when configured with a custom endpoint.
//
// To run:
//
//	go test -v -run ^TestHealthcheckWithCustomEndpoint$
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

// TestHealthcheckEndpoint verifies that the default healthcheck endpoint ("/healthcheck")
//
// responds with status 200 and body "OK" when the application is healthy.
//
// To run:
//
//	go test -v -run ^TestHealthcheckEndpoint$
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

// TestHealthcheckProbeFalse verifies that the healthcheck endpoint returns
//
//  HTTP 503 (Service Unavailable) when the Probe function returns false.
//
// To run:
//
//	go test -v -run ^TestHealthcheckProbeFalse$
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

// TestHealthcheckMethodNotAllowed verifies that the healthcheck endpoint
// 
// returns HTTP 405 (Method Not Allowed) when using a non-GET method.
//
// To run:
//
//	go test -v -run ^TestHealthcheckMethodNotAllowed$
func TestHealthcheckMethodNotAllowed(t *testing.T) {
	q := quick.New()
	q.Use(New(Options{
		App: q,
	}))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodPost, // POST instead of GET
		URI:    "/healthcheck",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := resp.AssertStatus(quick.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
}

// TestHealthcheckWithNextSkipping verifies that the healthcheck middleware
// 
// skips route logic when the Next function returns true, returning 404.
//
// To run:
//
//	go test -v -run ^TestHealthcheckWithNextSkipping$
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
