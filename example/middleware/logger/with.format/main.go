package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply the logger middleware with custom configuration
	q.Use(logger.New(logger.Config{
		Format:  "text", // Could it be "text", "json", "slog"
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG", // Could it be "DEBUG", "INFO", "WARN", "ERROR"
		CustomFields: map[string]string{
			"user_id": "12345",
			"trace":   "xyz",
		},
	}))

	// Define a GET route that logs request details
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Return a simple JSON response
		return c.Status(200).JSON(quick.M{
			"msg": "Quick ❤️",
		})
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XGET localhost:8080/v1/logger
