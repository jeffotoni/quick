package quick

import "testing"

// TestRouteHEAD tests whether a HEAD route returns the expected response headers and status
// without returning any response body.
//
// Usage:
//
//	go test -v -run TestRouteHEAD
func TestRouteHEAD(t *testing.T) {
	q := New()

	q.Head("/v1/user", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.String("Hello, HEAD!") // This should NOT be included in the body for HEAD
	})

	res, err := q.Qtest(QuickTestOptions{
		Method:  MethodHead,
		URI:     "/v1/user",
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	// The body for a HEAD request must be empty
	if res.BodyStr() != "" {
		t.Errorf("Expected empty body for HEAD, got '%s'", res.BodyStr())
	}

	// You can also check if the header is present
	if err := res.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Errorf("Expected Content-Type 'application/json', got '%v'", err)
	}
}
