package main

import (
	"context"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

// Example HTTP request using Quick and an external API (reqres.in).
// This server listens on port 3000 and forwards GET requests to an external API.
// - WithTimeout: Sets the HTTP client timeout to 10 seconds.
// - WithContext: Attaches a background context to the client which is not cancellable.
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
// - WithMaxConnsPerHost: Sets the maximum connections allowed per host (20).
// - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).

func main() {
	q := quick.New()

	// Define a GET endpoint that forwards requests to an external API.
	q.Get("/api/users", func(c *quick.Ctx) error {
		// Create an HTTP client with specific configurations.
		cClient := client.New(
			// Set the timeout for the HTTP client to 10 seconds.
			client.WithTimeout(10*time.Second),
			client.WithMaxConnsPerHost(20),
			client.WithDisableKeepAlives(false),
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

		// Perform a GET request to the external API.
		resp, err := cClient.Get("https://reqres.in/api/users/2")
		if err != nil {
			// Log the error and return a server error response if the GET request fails.
			log.Println("GET Error:", err)
			return c.Status(500).SendString("Failed to connect to external API")
		}

		// Log and return the response body from the external API.
		log.Println("GET Response:", string(resp.Body))
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start listening on port 3000 for incoming HTTP requests.
	q.Listen(":3000")
}

// curl --location 'http://localhost:3000/api/users' \
// --header 'Authorization: Bearer EXAMPLE_TOKEN'
