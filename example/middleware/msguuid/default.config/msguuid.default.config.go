package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/goquick/middleware/msguuid"
	"github.com/jeffotoni/goquick"
)

func main() {
	app := quick.New()

	app.Use(msguuid.New())

	app.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		fmt.Println("Headers:", c.Response.Header())
		return c.Status(200).JSON(nil)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}
