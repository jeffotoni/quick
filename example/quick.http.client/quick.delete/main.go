package main

import (
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

// Example HTTP request using Quick and an external API (reqres.in).
// This server listens on port 3000 and handles DELETE requests to delete user data on an external API.
// - WithTimeout: Sets the HTTP client timeout to 2 seconds.
// - WithHeaders: Adds custom headers for content type and authorization.

func main() {
	q := quick.New()

	// Define a DELETE endpoint to delete user data.
	q.Delete("/api/users/2", func(c *quick.Ctx) error {
		// Create an HTTP client with specific configurations.
		cClient := client.New(
			client.WithTimeout(2*time.Second),
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
		)

		// Perform a DELETE request to the external API.
		resp, err := cClient.Delete("https://reqres.in/api/users/2")
		if err != nil {
			log.Println("DELETE Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("DELETE Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}

// curl -X DELETE https://reqres.in/api/users/2 \
//       -H "Authorization: Bearer EXAMPLE_TOKEN"
