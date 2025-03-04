package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initializes the Quick framework

	// Sets a GET route for "/v1/user"
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")                       // Sets the content type as JSON
		return c.Status(200).SendString("Quick in action with Cors❤️!") // Returns a success message
	})

	// Sets a GET route for "/v2"
	q.Get("/v2", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")         // Sets the content type as JSON
		return c.Status(200).SendString("Is in the air!") // Returns a message indicating that the service is active
	})

	// Sets a GET route for "/v3"
	q.Get("/v3", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")   // Sets the content type as JSON
		return c.Status(200).SendString("Running!") // Returns a message confirming that the server is running
	})

	// Starts the server on port 8080, allowing connections from any IP
	q.Listen("0.0.0.0.0.0:8080")
}
