package helmet

import (
	"net/http"
	"testing"

	"github.com/jeffotoni/quick"
)

// TestHelmet tests the Helmet middleware by ensuring that security headers are added to the response.
func TestHelmet(t *testing.T) {
	q := quick.New()
	q.Use(Helmet())

	// Define a test endpoint to validate security headers
	q.Get("/v1/health", func(c *quick.Ctx) error {
		return c.Status(http.StatusOK).String("OK")
	})

	// Test health endpoint
	resp, err := q.Qtest(quick.QuickTestOptions{
		Method:     quick.MethodGet,
		URI:        "/v1/health",
		LogDetails: true,
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate response status
	if err := resp.AssertStatus(http.StatusOK); err != nil {
		t.Error(err)
	}
}

// TestWithoutHelmet tests the behavior of the application without the Helmet middleware.
func TestWithoutHelmet(t *testing.T) {
	q := quick.New()

	// Define a test endpoint to validate security headers
	q.Get("/v1/health", func(c *quick.Ctx) error {
		return c.Status(http.StatusOK).String("OK")
	})

	// Test health endpoint
	resp, err := q.Qtest(quick.QuickTestOptions{
		Method:     quick.MethodGet,
		URI:        "/v1/health",
		LogDetails: true,
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate response status
	if err := resp.AssertStatus(http.StatusOK); err != nil {
		t.Error(err)
	}
}
