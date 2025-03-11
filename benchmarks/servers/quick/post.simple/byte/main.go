package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initialize Quick framework

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {

		data := c.Body()

		return c.Status(200).Send(data)
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}
