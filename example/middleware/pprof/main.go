package main

import (
	"net/http"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	q := quick.New()

	// Apply the profiling middleware
	q.Use(pprof.New(
		pprof.Options{
			App: q,
		},
	))

	// Define a test route
	q.Get("/", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(http.StatusOK).JSON(map[string]any{
			"message": "Hello, World!",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
