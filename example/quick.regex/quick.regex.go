package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
	q := quick.New()

	// Adding middleware msgid
	q.Use(msgid.New())

	// Corrected route using :id instead of {id:[0-9]+}
	q.Get("/v1/user/:id", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).String("Quick ação total!!!")
	})

	q.Listen(":8080")
}
