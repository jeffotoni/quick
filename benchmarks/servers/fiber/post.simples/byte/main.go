package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/v1/user", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")

		rawBody := c.Body()
		return c.Send(rawBody)

	})

	app.Listen(":8080")
}
