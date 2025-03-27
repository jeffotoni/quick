package main

import (
	"net/http"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	q := quick.New()

	// Apply the profiling middleware
	q.Use(pprof.New())

	// Define a test route
	q.Get("/debug/pprof*", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(http.StatusOK).JSON(quick.M{
			"message": "Hello, World!",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
