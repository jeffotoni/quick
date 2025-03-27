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
			return c.Status(quick.StatusTooManyRequests).JSON(response)
		},
	}))

	// Simulate multiple requests to test rate limiting
	for i := 0; i < 5; i++ { // Enviar mais requisições que o limite (3)
		res, _ := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/",
		})
		fmt.Println(res.BodyStr())
	}

	// Out put:
	// {"msg":"Quick in action ❤️!"}
	// {"msg":"Quick in action ❤️!"}
	// {"msg":"Quick in action ❤️!"}
	// Rate Limit Exceeded: map[error:Too many requests message:You have exceeded the request limit. Please wait 10 seconds and try again. retry_after:10s]
	// {"error":"Too many requests","message":"You have exceeded the request limit. Please wait 10 seconds and try again.","retry_after":"10s"}
}

// This function is named ExampleRateLimiter_group()
// it with the Examples type.
func ExampleRateLimiter_group() {
	// Create a new Quick instance
	q := quick.New()

	// Rate Limiter Middleware
	limiterMiddleware := New(Config{
		// Maximum 5 requests allowed per IP address within a 10-second window
		Max: 5,
		// The limit resets every 10 seconds
		Expiration: 10 * time.Second,
		// Use the client's IP address as the unique key to track rate limits
		KeyGenerator: func(c *quick.Ctx) string {
			return c.RemoteIP()
		},
		// If the rate limit is exceeded, send an error message and instructions
		LimitReached: func(c *quick.Ctx) error {
			// Set content type to JSON
			c.Set("Content-Type", "application/json")
			// Set the Retry-After header to indicate how long the client should wait before retrying
			c.Set("Retry-After", "10") // The client should wait 10 seconds before retrying
			// Response structure
			response := map[string]string{
				"error":       "Too many requests",
				"message":     "You have exceeded the request limit. Please wait 10 seconds and try again.",
				"retry_after": "10s",
			}

			// Log to verify that the rate limit exceeded response is being sent
			log.Println("Rate Limit Exceeded:", response)

			// Return the response with HTTP status 429 (Too Many Requests)
			return c.Status(quick.StatusTooManyRequests).JSON(response)
		},
	})

	// Create an API group with rate limit middleware
	api := q.Group("/v1")
	// Apply the rate limiter middleware to the /api group
	api.Use(limiterMiddleware)

	// Define route /api/users that responds with a list of users
	api.Get("/users", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "list of users"})
	})

	// Define route without rate limit
	// This route is not affected by the rate limiter
	q.Get("/", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "Quick in action ❤️!"})
	})

	// Functionality test using QuickTest (simulates a request for the protected route)
	for i := 0; i < 5; i++ { // Send more requests than the limit (3) to test the block
		res, _ := q.Qtest(quick.QuickTestOptions{
			Method: quick.MethodGet,
			URI:    "/v1/users",
		})
		fmt.Println(res.BodyStr())
	}

	// Out put:
	// {"msg":"list of users"}
	// {"msg":"list of users"}
	// {"msg":"list of users"}
	// Rate Limit Exceeded: map[error:Too many requests message:You have exceeded the request limit. Please wait 10 seconds and try again. retry_after:10s]
	// {"error":"Too many requests","message":"You have exceeded the request limit. Please wait 10 seconds and try again.","retry_after":"10s"}
}
