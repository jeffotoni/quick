package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New(quick.Config{
		MaxBodySize: 5 * 1024 * 1024,
	})

	group := q.Group("/v1")
	group.Get("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
	})

	group.Post("/user", func(c *quick.Ctx) error {
		return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
	})

	group2 := q.Group("/v2")

	group2.Get("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action com [GET] /v2/user!")
	})

	group2.Post("/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Quick in action com [POST] /v2/user!")
	})

	q.Listen("0.0.0.0:8080")
}
