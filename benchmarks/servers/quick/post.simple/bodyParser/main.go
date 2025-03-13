package main

import "github.com/jeffotoni/quick"

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

// $ curl --location 'http://localhost:8080/v1/user' \
// --header 'Content-Type: application/json' \
// --data '{"name": "Alice", "year": 20}'
func main() {
	q := quick.New() // Initialize Quick framework

	// Define a POST route at /v1/user
	q.Post("/v1/user", func(c *quick.Ctx) error {
		var my My // Create a variable to store incoming user data

		// Parse the request body into the struct
		err := c.BodyParser(&my)
		if err != nil {
			// If parsing fails, return a 400 Bad Request response
			return c.Status(400).SendString(err.Error())
		}

		// Return the parsed JSON data as a response with 200 OK
		return c.Status(200).JSON(my)

		// Alternative:
		// return c.Status(200).String(c.BodyString())
		// Return raw request body as a string
	})

	// Start the server and listen on port 8080
	q.Listen("0.0.0.0:8080")
}
