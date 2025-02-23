package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Static("/static", "./static")

	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/*")
		return nil
	})

	q.Listen("0.0.0.0:8080")
}
