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
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
// - WithHeaders: Adds custom headers (e.g., Authorization: "Bearer token").
func main() {
	// Create a new Quick HTTP client with retry and custom headers
	cC := client.New(
		// Enables automatic retries on failed requests
		client.WithRetry(client.RetryConfig{
			MaxRetries: 3,                         // Maximum number of retry attempts
			Delay:      1 * time.Second,           // Delay before each retry
			UseBackoff: false,                     // Disables exponential backoff (fixed retry delay)
			Statuses:   []int{502, 503, 504, 403}, // Retries only on these HTTP status codes
			FailoverURLs: []string{ // Backup URLs to try if the primary request fails
				"http://backup1",
				"https://reqres.in/api/users",
				"https://httpbin.org/post",
			},
			EnableLog: true, // Enables logging for debugging retries
		}),
		// Sets custom headers for the request
		client.WithHeaders(map[string]string{
			"Authorization": "Bearer token", // Adds an authentication token
		}),
	)

	// Perform the POST request
	resp, err := cC.Post("https://httpbin_error.org/post", map[string]string{
		"name":  "Jefferson",
		"email": "jeff@example.com",
	})
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}

	// Print the response body and status code
	fmt.Println("POST Response Status:", resp.StatusCode)
	fmt.Println("POST Response Body:", string(resp.Body))

}
