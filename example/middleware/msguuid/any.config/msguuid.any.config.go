package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/msguuid"
)

func main() {
	app := quick.New()

	// Enable MsgUUID Middleware
	app.Use(msguuid.New())

	// Define an endpoint that includes MsgUUID
	app.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log the response headers to check the UUID
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status with no body
		return c.Status(200).JSON(nil)
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}

//curl --location 'http://localhost:8080/v1/user'
