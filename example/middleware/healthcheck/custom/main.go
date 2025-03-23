package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/healthcheck"
)

func main() {
	q := quick.New()

	// register healthcheck middleware with custom endpoint
	q.Use(healthcheck.New(
		healthcheck.Options{
			Endpoint: "/health",
			Probe: func(c *quick.Ctx) bool {
				return true
			},
			App: q,
		},
	))

	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Home page")
	})

	log.Fatalln(q.Listen(":8080"))
}
