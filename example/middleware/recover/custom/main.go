package main

import (
	"errors"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/recover"
)

func main() {
	q := quick.New()

	// Apply the Recover middleware with custom configuration
	q.Use(recover.New(recover.Config{
		Next: func(c *quick.Ctx) bool {
			return false
		},
		EnableStacktrace: true,
	}))

	// Define a test route
	q.Get("/v1/panic", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// halt the server
		panic(errors.New("Panicking!"))
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
