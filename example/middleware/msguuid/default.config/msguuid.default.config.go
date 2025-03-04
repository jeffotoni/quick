package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	app := quick.New()

	// Apply MsgUUID Middleware globally
	app.Use(msguuid.New())

	// Define an endpoint that responds with a UUID
	app.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log headers to validate UUID presence
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status
		return c.Status(200).JSON(nil)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}
