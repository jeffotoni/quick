package logger

import "github.com/jeffotoni/quick"

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {

	q := quick.New()

	// Apply the logger middleware with default configuration
	q.Use(New())

	// Define a simple GET route
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "Quick ❤️",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
