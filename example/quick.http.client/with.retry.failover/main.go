package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
//   - WithTimeout: Sets the HTTP client timeout to 30 seconds.
//   - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
//   - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
func main() {
	// Create a new HTTP client with specific configurations.
	cClient := client.New(
		// Set a timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, 'Content-Type' is set to 'application/json' which is typical for API calls.
		client.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),

		// Enable automatic retry mechanism with specific configurations.
		// This is useful for handling intermittent errors and ensuring robustness.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,                    // Retry failed requests up to two times.
				Delay:      1 * time.Second,      // Wait for 1 second before retrying.
				UseBackoff: true,                 // Use exponential backoff strategy for retries.
				Statuses:   []int{500, 502, 503}, // HTTP status codes that trigger a retry.
				FailoverURLs: []string{ // Alternate URLs to try if the main request fails.
					"https://httpbin_error.org/post",
					"https://httpbin.org/post",
				},
				EnableLog: true, // Enable logging for retry operations.
			}),
	)

	// Perform a POST request using the configured HTTP client.
	// Includes a JSON payload with a "name" key.
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{
		"name": "jeffotoni in action with Quick!!!",
	})

	// Check if there was an error with the POST request.
	if err != nil {
		// If an error occurs, log the error and terminate the program.
		log.Fatalf("POST request failed: %v", err)
	}

	// Print the response from the server to the console.
	fmt.Println("POST Form Response:", string(resp.Body))
}
