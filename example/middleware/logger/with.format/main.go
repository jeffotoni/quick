package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

// curl -i -XGET localhost:8080/v1/logger
func main() {

	q := quick.New()

	q.Use(logger.New(logger.Config{
		Format:  "text", // Could it be "text", "json", "slog"
		Pattern: "[${level}] [${time}] ${ip} ${method} ${path} ${status} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG", // Could it be "DEBUG", "INFO", "WARN", "ERROR"
		CustomFields: map[string]string{
			"user_id": "12345",
			"trace":   "xyz",
		},
	}))

	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type my struct {
			Msg string `json:"msg"`
		}

		return c.Status(200).JSON(&my{
			Msg: "Quick ❤️",
		})
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}
