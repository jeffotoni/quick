package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

// curl -i -XGET localhost:8080/v1/logger
func main() {

	q := quick.New()
	q.Use(logger.New())

	q.Use(logger.New(logger.Config{
		Level: "DEGUB",
	}))

	q.Use(logger.New(logger.Config{
		Level: "WARN",
	}))

	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "Quick",
		})
	})

	q.Listen("0.0.0.0:8080")
}
