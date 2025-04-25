package limiter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jeffotoni/quick"
)
// TestLimiterMiddleware verifies the complete behavior of the rate limiting middleware.
//
// The test validates:
//   - Request allowance within the defined limit
//   - Proper blocking of requests exceeding the limit
//   - Correct HTTP status codes (200 vs 429)
//   - Rate limit window expiration and reset
//   - Integration with the Quick framework
//
// Test Flow:
// 1. Configures middleware with 3 requests per 2 seconds limit
// 2. Makes 6 requests in sequence
// 3. Verifies first 3 requests succeed (200 OK)
// 4. Verifies next 3 requests are blocked (429 Too Many Requests)
// 5. Waits for rate limit window to reset (2 seconds)
// 6. Verifies requests are allowed again after reset
//
// Debugging:
//   - The test includes debug logs to trace middleware execution
//   - Key generation and limit reached callbacks are instrumented
func TestLimiterMiddleware(t *testing.T) {
	// Create a new Quick instance
	q := quick.New()

	// Configure rate limiting middleware with:
	// - Maximum 3 requests per window
	// - 2 second window duration
	// - Client IP as rate limit key
	// - Custom rate limit response
	q.Use(New(Config{
		Max:        3,               // Allow up to 3 requests
		Expiration: 2 * time.Second, // Reset after 2 seconds
		KeyGenerator: func(c *quick.Ctx) string {
			fmt.Println("I'm here KeyGenerator........:", c.RemoteIP())
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) error {
			fmt.Println("I'm here LimitReached........")
			return c.Status(http.StatusTooManyRequests).SendString(`{"error":"Too many requests"}`)
		},
	}))

	// Register a test route
	q.Get("/", func(c *quick.Ctx) error {
		t.Log("[DEBUG] Handler executed") // Log handler execution
		return c.Status(200).JSON(map[string]string{"msg": "Hello, Quick!"})
	})

	// Start test server
	ts := httptest.NewServer(q)
	defer ts.Close()

	client := ts.Client()

	for i := 0; i < 6; i++ { // We make 6 requests to ensure the limit of 3 is reached
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		if i >= 3 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("[DEBUG] Request %d: Expected 429, got %d", i+1, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Perform the blocked request (should be 429)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error on rate-limited request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got %d", resp.StatusCode)
	}

	// Wait for the expiration period
	t.Log("[DEBUG] Waiting for expiration period...")
	time.Sleep(2 * time.Second)

	// Perform a request after expiration (should be 200 again)
	resp, err = client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error after expiration period: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK after expiration, got %d", resp.StatusCode)
	}
}
