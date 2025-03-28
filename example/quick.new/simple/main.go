package main

import (
	"github.com/jeffotoni/quick"
)

// This example demonstrates simple usage of Quick
func main() {
	// Initialize a new Quick instance
	q := quick.New()

	// Define a simple GET route at the root path
	q.Get("/", func(c *quick.Ctx) error {
		// Set response header to indicate plain text response
		c.Set("Content-Type", "text/plain")

		// Return a 200 OK response with a message
		return c.Status(200).String("Quick in action!")
	})

	// Start the Quick server on port 8080
	q.Listen(":8080")
}

// $ curl -X GET http://localhost:8080/ -H "Content-Type: text/plain"
