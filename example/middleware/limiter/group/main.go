package main

import (
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	// Create a new Quick instance
	q := quick.New()

	// Rate Limiter Middleware
	limiterMiddleware := limiter.New(limiter.Config{
		// Maximum 3 requests allowed per IP address within a 10-second window
		Max: 3,
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

	// Define route /api/posts that responds with a list of posts
	api.Get("/posts", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "list of posts"})
	})

	// Define route without rate limit
	// This route is not affected by the rate limiter
	q.Get("/", func(c *quick.Ctx) error {
		return c.JSON(map[string]string{"msg": "Quick in action!"})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

// To test without Rate Limiter, use the following curl command:
// $ curl --location 'http://localhost:8080/'

// To test with Rate Limiter, use these curl commands:
// $ curl --location 'http://localhost:8080/v1/users'
// $ curl --location 'http://localhost:8080/v1/posts'

//Script
// async function testRateLimiter() {
//     function delay(ms) {
//         return new Promise(resolve => setTimeout(resolve, ms));
//     }

//     for (let i = 0; i < 10; i++) {
//         await pm.sendRequest("http://localhost:8080/v1/users", function (err, res) {
//             console.log(`Request ${i + 1}: Status - ${res.code} | Body - ${res.text()}`);
//         });
//         await delay(200); // Aumente o delay para 200ms para evitar que todas as requisições sejam tratadas ao mesmo tempo
//     }
// }
