package main

import (
	"net/http"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	// Create a new Quick application instance
	q := quick.New()

	// Apply the pprof middleware to enable runtime profiling
	// This allows access to profiling endpoints like /debug/pprof/heap, /goroutine, etc.
	q.Use(pprof.New())

	// Define a test route that matches /debug/pprof*
	// This is required so that the Quick router delegates the request to the pprof middleware
	q.Get("/debug/pprof*", func(c *quick.Ctx) error {
		// Set the response content type to JSON
		c.Set("Content-Type", "application/json")

		// Return a basic JSON response for testing
		return c.Status(http.StatusOK).JSON(quick.M{
			"message": "Hello, World!",
		})
	})

	// Start the server on all interfaces (localhost and external) at port 8080
	q.Listen("0.0.0.0:8080")
}

// $ curl http://localhost:8080/debug/pprof/
