package compress

import (
	"compress/gzip"
	"io"
	"testing"

	"github.com/jeffotoni/quick"
)

// TestGzipMiddleware validates the Gzip middleware functionality in Quick.
//
// This test suite ensures that the middleware properly compresses responses
// when the client supports gzip and leaves responses uncompressed when
// the client does not request gzip encoding.
//
// The test executes two scenarios:
// 1. A request with "Accept-Encoding: gzip" should return a gzipped response.
// 2. A request without "Accept-Encoding" should return a normal (uncompressed) response.
//
// Middleware Setup:
// The test initializes a new Quick instance, applies the Gzip middleware,
// and registers a test endpoint "/v1/user".
//
// Example Usage:
//   - Run `go test -v` to execute the test suite.
//
// Assertions:
//   - The response must have "Content-Encoding: gzip" when requested.
//   - The body should be successfully decompressed and match the expected output.
//   - The response remains uncompressed when no "Accept-Encoding" header is provided.
//
// Errors:
//   - If gzip encoding fails, the test fails with an appropriate error message.
//   - If the response body does not match the expected content, the test reports failure.
func TestGzipMiddleware(t *testing.T) {
	q := quick.New()
	q.Use(Gzip())

	// Registering a test route that returns a JSON response
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(quick.StatusOK).String("Hello, Quick!")
	})

	// Scenario 1: Response should be gzipped when "Accept-Encoding: gzip" is present
	t.Run("Response should be gzipped when Accept-Encoding includes gzip", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method:  quick.MethodGet,
			URI:     "/v1/user",
			Headers: map[string]string{"Accept-Encoding": "gzip"},
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Ensure response status is OK (200)
		if err := resp.AssertStatus(quick.StatusOK); err != nil {
			t.Error(err)
		}

		// Ensure response is gzipped
		if err := resp.AssertHeader("Content-Encoding", "gzip"); err != nil {
			t.Error(err)
		}

		// Attempt to decompress the gzipped response
		gzr, err := gzip.NewReader(resp.Response().Body)
		if err != nil {
			t.Errorf("Failed to create gzip reader: %v", err)
		}
		defer gzr.Close()

		// Read and validate the decompressed response body
		body, err := io.ReadAll(gzr)
		if err != nil {
			t.Errorf("Failed to read gzip body: %v", err)
		}
		if string(body) != "Hello, Quick!" {
			t.Errorf("Expected body 'Hello, Quick!', got '%s'", string(body))
		}
	})

	// Scenario 2: Response should not be gzipped when "Accept-Encoding" is missing
	t.Run("Response should not be gzipped when Accept-Encoding is missing", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/v1/user",
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Ensure response status is OK (200)
		if err := resp.AssertStatus(quick.StatusOK); err != nil {
			t.Error(err)
		}

		// Ensure response does not have gzip encoding
		if err := resp.AssertHeader("Content-Encoding", ""); err != nil {
			t.Error(err)
		}

		// Ensure the response body matches the expected uncompressed output
		if err := resp.AssertBodyContains("Hello, Quick!"); err != nil {
			t.Error(err)
		}
	})
}
