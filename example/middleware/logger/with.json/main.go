package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply logger with JSON format
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "INFO",
	}))

	q.Use(logger.New(logger.Config{
		Format:  "json",
		Pattern: "[${level}] ${time} ${ip} ${method} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-debug",
		},
	}))

	// Apply the logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "WARN",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "usr-001",
			"trace":   "trace-warn",
		},
	}))

	// Define an endpoint that triggers logging
	q.Get("/v1/logger/json", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "JSON logging example",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XGET localhost:8080/v1/logger/json
