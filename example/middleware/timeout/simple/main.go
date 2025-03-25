package main

import (
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/timeout"
)

func main() {
	q := quick.New()

	// Apply the Timeout middleware
	q.Use(timeout.New(timeout.Options{
		Duration: 2 * time.Second,
	}))

	// Define a test route
	q.Get("/v1/slow", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Simulate a slow response
		time.Sleep(5 * time.Second)
		return c.SendString("Slow response")
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
