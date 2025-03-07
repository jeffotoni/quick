package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

// Example HTTP request using Quick and an external API (reqres.in).
// This server listens on port 3000 and handles PUT requests to update user data on an external API.
// - WithTimeout: Sets the HTTP client timeout to 10 seconds.
// - WithContext: Attaches a background context to the client which is not cancellable.
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
// - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).

func main() {
	q := quick.New()

	// Define a PUT endpoint to update user data.
	q.Put("/api/users/2", func(c *quick.Ctx) error {
		// Read the request body from the client
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Read Error:", err)
			return c.Status(500).SendString("Failed to read request body")
		}

		// Create an HTTP client with specific configurations.
		cClient := client.New(
			// Set the timeout for the HTTP client to 10 seconds.
			client.WithTimeout(10*time.Second),
			// Add custom headers, including content type and authorization token.
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
			// Use a background context for the HTTP client. This context cannot be cancelled
			// and does not carry any deadline. It is suitable for operations that run
			// indefinitely or until the application is shut down.
			client.WithContext(context.Background()),
		)

		// Perform a PUT request to the external API with the data received from the client.
		resp, err := cClient.Put("https://reqres.in/api/users/2", requestBody)
		if err != nil {
			log.Println("PUT Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("PUT Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}

// curl -X PUT https://reqres.in/api/users/2 \
//      -H "Content-Type: application/json" \
//      -d '{"name": "Morpheus", "job": "zion resident"}'
