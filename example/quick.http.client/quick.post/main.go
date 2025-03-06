package main

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

// Example HTTP request using Quick and an external API (reqres.in).
// This server listens on port 3000 and forwards POST requests to an external API.
// - WithTimeout: Sets the HTTP client timeout to 5 seconds.
// - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
func main() {
	// Initialize the Quick framework.
	q := quick.New()

	// Define a POST endpoint to process incoming requests.
	q.Post("/api/users", func(c *quick.Ctx) error {
		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString("Error reading request body: " + err.Error())
		}

		// Check if the request body is empty
		if len(body) == 0 {
			return c.Status(400).SendString("Error: Request body is empty")
		}

		// Validate that the request body is valid JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			return c.Status(400).SendString("Error: Invalid JSON")
		}

		// Create a modular HTTP client with customizable options.
		cClient := client.New(
			// Sets the HTTP timeout to 5 seconds.
			client.WithTimeout(5*time.Second),

			// Enables or disables HTTP Keep-Alive connections (false = keep-alives enabled).
			client.WithDisableKeepAlives(false),

			// Adds custom headers to the request, including Content-Type and Authorization.
			client.WithHeaders(map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer EXAMPLE_TOKEN",
			}),
		)

		// Forward the request to the external API
		resp, err := cClient.Post("https://reqres.in/api/users", json.RawMessage(body))
		if err != nil {
			log.Println("Error making request to external API:", err)
			return c.Status(500).SendString("Error connecting to external API")
		}

		// Log response from external API for debugging
		log.Println("External API response:", string(resp.Body))

		// Return the response from the external API to the client
		return c.Status(resp.StatusCode).Send(resp.Body)
	})

	// Start the server on port 3000
	q.Listen(":3000")
}

// curl -X POST http://localhost:3000/api/users \
//      -H "Content-Type: application/json" \
//      -H "Authorization: Bearer EXAMPLE_TOKEN" \
//      -d '{"name": "John Doe", "job": "Software Engineer"}'
