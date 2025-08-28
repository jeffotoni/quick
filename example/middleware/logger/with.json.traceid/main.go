package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

var NAME_TRACE_ID = "X-TRACE-ID"

func main() {

	q := quick.New()

	// Apply logger with JSON format
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "DEBUG",
		// Pattern: "[${X-TRACE-ID} ${service} [${time}] ${level} ${method} ${path} ${status} - ${latency}\n", // text and slog
	}))

	q.Post("/v1/logger/json", func(c *quick.Ctx) error {
		defer c.SaveContext() // ← MÁGICA! Salva todos os dados no final

		c.Set("Content-Type", "application/json")

		traceID := c.GetTraceID(NAME_TRACE_ID)

		// set header and context in request
		c.SetTraceID(NAME_TRACE_ID, traceID)

		// request - agora acumula localmente sem notificar logger ainda
		c.SetContext().
			Str("service-1", "user-service-1").
			Str("func-error-1", "error func2")

		c.SetContext().
			Str("service-2", "user-service-2").
			Str("func-error-2", "error func2").
			Int("service-int", 3939).
			Bool("service-bool", true)

		return c.Status(200).JSON(quick.M{
			"msg":     "JSON logging example",
			"traceID": traceID,
			"context": fmt.Sprintf("%v", c.GetAllContextData()),
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

// $ curl -i -XPOST localhost:8080/v1/logger/json id '{"name": "@jeffotoni"}'
