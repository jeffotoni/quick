package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

var NAME_TRACE_ID = "X-TRACE-ID"

func main() {

	q := quick.New()

	// Apply logger with JSON format
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "INFO",
	}))

	// Define an endpoint that triggers logging
	q.Post("/v1/logger/json", func(c *quick.Ctx) error {

		c.Set("Content-Type", "application/json")

		traceID := c.GetTraceID(NAME_TRACE_ID)

		// request
		contextData := map[string]string{
			NAME_TRACE_ID: traceID,
			"service":     "user-service",
			"function":    "createUser",
		}

		// set request
		c.SetContext(contextData)

		// response
		c.Set(NAME_TRACE_ID, traceID)

		return c.Status(200).JSON(quick.M{
			"msg":     "JSON logging example",
			"traceID": traceID,
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XPOST localhost:8080/v1/logger/json id '{"name": "@jeffotoni"}'
