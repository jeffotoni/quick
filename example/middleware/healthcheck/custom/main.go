package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/healthcheck"
)

func main() {
	app := quick.New()

	app.Use(healthcheck.New(healthcheck.Options{
		App:      app,
		Endpoint: "/health",
		Probe: func(c *quick.Ctx) bool {
			return true
		},
	}))

	app.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Home page")
	})

	log.Fatal(app.Listen(":8080"))
}
