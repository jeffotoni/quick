package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New()
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action com Cors❤️!")
	})

	q.Get("/v2", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Está no ar️!")
	})

	q.Get("/v3", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Rodando️!")
	})

	q.Listen("0.0.0.0:8080")
}
