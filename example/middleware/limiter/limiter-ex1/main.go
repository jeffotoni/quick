package main

import (
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	q := quick.New()

	// Apply the rate limiter middleware
	q.Use(limiter.New(limiter.Config{
		// Maximum 5 requests allowed per IP
		Max: 10,
		// The limit resets after 5 seconds
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *quick.Ctx) string {
			// Uses the client's IP address as the key
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) error {
			c.Set("Content-Type", "application/json")
			c.Set("Retry-After", "5") // The client should wait 5 seconds before retrying
			return c.Status(quick.StatusTooManyRequests).JSON(map[string]string{
				"error":       "Too many requests",
				"message":     "You have exceeded the request limit. Please wait 1 second and try again.",
				"retry_after": "1s", // Suggests a 1-second delay before retrying
			})
		},
	}))

	// Define a simple GET route
	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Hello, Quick!"})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

// $ curl --location 'http://localhost:8080/'

//Here is the direct curl script to test the Rate Limiter:

// Function to test the Rate Limiter
/*async function testRateLimiter() {

    // Helper function to create a delay (pause execution for a given time)
    function delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    // Loop to send 10 requests to the API
    for (let i = 0; i < 10; i++) {
        pm.sendRequest("http://localhost:8080/", function (err, res) {
            // Logs the request number, HTTP status code, and response body
            console.log(`Request ${i + 1}: Status - ${res.code} | Body - ${res.text()}`);
        });

        await delay(200); // Adds a 200ms delay to prevent sending all requests at once
    }
}

// Calls the function to start the test
testRateLimiter();
*/
