package main

import (
	"log"

	"github.com/jeffotoni/goquick/middleware/logger"
	"github.com/jeffotoni/quick"
)

// curl -i -XGET localhost:8080/v1/logger
func main() {

	q := quick.New()
	q.Use(logger.New())

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
