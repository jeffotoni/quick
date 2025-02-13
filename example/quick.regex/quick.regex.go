package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
	q := quick.New()

	// adicionando middleware msgid
	q.Use(msgid.New())

	q.Get("/v1/user/{id:[0-9]+}", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).String("Quick ação total!!!")
	})

	q.Use(msgid.New())

	q.Get("/v2/tipos/{id:[0-9]+}", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).String("Quick funcionando!!!")
	})

	q.Listen(":8080")
}
