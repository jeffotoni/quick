package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

func main() {

	q := quick.New(
		quick.Config{
			NoBanner: false,
		}) // Initialize Quick framework

	// Route to greet a user by name (dynamic route parameter)
	q.Get("/v1/user/:name", func(c *quick.Ctx) error {
		name := c.Param("name")                              // Retrieve the 'name' parameter from the URL
		c.Set("Content-Type", "text/plain")                  // Set response content type as plain text
		return c.Status(200).SendString("Olá " + name + "!") // Return greeting message
	})

	// Simple route returning a static message
	q.Get("/v2/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")            // Set response content type as JSON
		return c.Status(200).SendString("Opa, funcionando!") // Return confirmation message
	})

	// Route to return an ID from the URL
	q.Get("/v3/user/:id", func(c *quick.Ctx) error {
		id := c.Param("id")                         // Retrieve the 'id' parameter from the URL
		c.Set("Content-Type", "application/json")   // Set response content type as JSON
		return c.Status(200).SendString("Id:" + id) // Return the ID in the response
	})

	// Complex route with multiple parameters
	q.Get("/v1/userx/:p1/:p2/cust/:p3/:p4", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")           // Set response content type as JSON
		return c.Status(200).SendString("Quick in action!") // Return a success message
	})

	// Print all registered routes
	for k, v := range q.GetRoute() {
		fmt.Println(k, "[", v, "]")
	}

	// Start the server and listen on port 8080
	q.Listen(":8394")
}
