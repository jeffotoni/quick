package main

import (
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

// User struct defines a user with Name and Year of birth
type User struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	q := quick.New()

	// Simulating a "database" with pre-registered users
	users := map[string]User{
		"1": {Name: "Maria", Year: 2000}, // Fixed user with ID 1
	}

	// DELETE route to remove a user by ID
	q.Delete("/v1/user/:id", func(c *quick.Ctx) error {
		userID := c.Params["id"] // Retrieve user ID from URL parameter

		// Check if the user exists in the "database"
		if _, exists := users[userID]; !exists {
			return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "User not found"})
		}

		// Delete the user from the "database"
		delete(users, userID)

		// Return a success response
		return c.Status(http.StatusOK).JSON(map[string]string{"msg": "User deleted successfully!"})
	})

	// Start the server on port 8080
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
