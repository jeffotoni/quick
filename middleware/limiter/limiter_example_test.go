package limiter

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleRateLimiter()
// it with the Examples type.
func ExampleRateLimiter() {
	// Initialize the Quick framework
	q := quick.New()

	// Apply the rate limiter middleware with max 3 requests per IP address within 10 seconds
	q.Use(New(Config{
		Max:        3,                // Max 3 requests per IP
		Expiration: 10 * time.Second, // Reset the rate limit every 10 seconds
		KeyGenerator: func(c *quick.Ctx) string {
			return c.RemoteIP() // Use the client's IP address as the key
		},
		LimitReached: func(c *quick.Ctx) error {
			c.Set("Content-Type", "application/json")
			c.Set("Retry-After", "10") // Instruct client to wait 10 seconds before retrying
			response := map[string]string{
				"error":       "Too many requests",
				"message":     "You have exceeded the request limit. Please wait 10 seconds and try again.",
				"retry_after": "10s",
			}
			// Log for verification
			log.Println("Rate Limit Exceeded:", response)

			return c.Status(quick.StatusTooManyRequests).JSON(response)
		},
	}))

	// Create a simple GET route with rate limiting
	q.Get("/", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "Quick in action ❤️!"})
	})

	// Test the rate limiting using QuickTest (Simulate a request to the defined route for testing)
	res, _ := q.QuickTest("GET", "/", nil)
	fmt.Println(res.BodyStr())

	// Out put:
	// Rate Limit Exceeded: map[error:Too many requests message:You have exceeded the request limit.
	// Please wait 10 seconds and try again. retry_after:10s]

}
