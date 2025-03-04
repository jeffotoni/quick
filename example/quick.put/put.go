package main

import "github.com/jeffotoni/quick"

func main() {
	q := quick.New() // Initialize Quick framework

	// PUT route to update a user by ID
	q.Put("/users/:id", func(c *quick.Ctx) error {
		userID := c.Param("id") // Retrieve the user ID from the URL parameter
		// Logic to update user data would go here
		return c.Status(200).SendString("User " + userID + " updated successfully!")
	})

	// PUT route to update a specific type by ID
	q.Put("/tipos/:id", func(c *quick.Ctx) error {
		tiposID := c.Param("id") // Retrieve the type ID from the URL parameter
		// Logic to update the type would go here
		return c.Status(200).SendString("User " + tiposID + " type updated successfully!")
	})

	// Start the server and listen on port 8080
	q.Listen(":8080")
}
