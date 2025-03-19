package maxbody

import (
	"io"
	"net/http"
	"testing"

	"github.com/jeffotoni/quick"
)

const defaultMaxBytesTest int64 = 10 // 10 bytes

// TestBodySizeMiddleware validates the middleware's ability to enforce request body size limits.
//
// This test suite evaluates various scenarios, including requests:
// 1. Within the allowed size limit.
// 2. Exceeding the allowed size limit.
// 3. Exactly at the allowed size limit.
// 4. With an empty body.
//
// Each test ensures the middleware correctly handles HTTP request bodies and returns the expected status codes.
func TestBodySizeMiddleware(t *testing.T) {
	q := quick.New()
	q.Use(New(defaultMaxBytesTest))

	// Define a test endpoint to validate request handling
	q.Post("/v1/upload", func(c *quick.Ctx) error {
		// Ensure the request body is fully read or discarded to trigger middleware validation
		_, _ = io.Copy(io.Discard, c.Request.Body)

		c.Set("Content-Type", "text/plain")
		return c.Status(http.StatusOK).String("Upload successful")
	})

	// Scenario 1: Request within the allowed size limit
	// Expected: HTTP 200 OK
	t.Run("Allow request within limit", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method:     quick.MethodPost,
			URI:        "/v1/upload",
			Body:       []byte("123456"), // 6 bytes, within the 10-byte limit
			LogDetails: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := resp.AssertStatus(http.StatusOK); err != nil {
			t.Error(err)
		}
	})

	// Scenario 2: Request exceeding the allowed size limit
	// Expected: HTTP 413 Request Entity Too Large
	t.Run("Reject request exceeding limit", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodPost,
			URI:    "/v1/upload",
			Body:   []byte("123456789023"), // 12 bytes, exceeding the 10-byte limit
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := resp.AssertStatus(http.StatusRequestEntityTooLarge); err != nil {
			t.Error(err)
		}
	})

	// Scenario 3: Request exactly at the allowed size limit
	// Expected: HTTP 200 OK
	t.Run("Allow request at exact limit", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method:     quick.MethodPost,
			URI:        "/v1/upload",
			Body:       []byte("1234567890"), // Exactly 10 bytes, matching the limit
			LogDetails: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Expecting HTTP 200 OK since the request body is exactly at the limit
		if err := resp.AssertStatus(http.StatusOK); err != nil {
			t.Error(err)
		}
	})

	// Scenario 4: Request with no body
	// Expected: HTTP 200 OK (Empty body should be accepted)
	t.Run("Allow request with no body", func(t *testing.T) {
		resp, err := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodPost,
			URI:    "/v1/upload",
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := resp.AssertStatus(http.StatusOK); err != nil {
			t.Error(err)
		}
	})
}
