package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/healthcheck"
)

func main() {
	q := quick.New()

	// register healthcheck middleware
	q.Use(healthcheck.New(
		healthcheck.Options{
			App: q,
		},
	))

	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Home page")
	})

	log.Fatalln(q.Listen(":8080"))
}
