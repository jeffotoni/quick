package healthcheck

import (
	"testing"

	"github.com/jeffotoni/quick"
)

// TestHealthcheck tests the healthcheck middleware with custom endpoint.
func TestHealthcheckWithCustomEndpoint(t *testing.T) {
	q := quick.New()
	q.Use(New(
		Options{
			Endpoint: "/v1/health",
			App:      q,
			Probe: func(c *quick.Ctx) bool {
				return true
			},
		},
	))

	resp, err := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/health",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.StatusCode() != quick.StatusOK {
		t.Errorf("Expected status code %d, got %d", quick.StatusOK, resp.StatusCode())
	}

	if string(resp.Body()) != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", string(resp.Body()))
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

	if resp.StatusCode() != quick.StatusOK {
		t.Errorf("Expected status code %d, got %d", quick.StatusOK, resp.StatusCode())
	}

	if string(resp.Body()) != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", string(resp.Body()))
	}
}
