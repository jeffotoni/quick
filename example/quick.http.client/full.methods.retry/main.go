package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
//   - WithTimeout: Sets the HTTP client timeout to 30 seconds.
//   - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).
//   - WithMaxIdleConns: Defines the maximum number of idle connections (20).
//   - WithMaxConnsPerHost: Sets the maximum connections allowed per host (20).
//   - WithMaxIdleConnsPerHost: Sets the maximum number of idle connections per host (20).
//   - WithContext: Injects a context for the client (context.TODO() used as placeholder).
//   - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
//   - WithTLSConfig: Configures TLS settings, including InsecureSkipVerify and TLS version.
//   - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
//     with exponential backoff (2s-bex) and a maximum of 3 attempts.
func main() {

	// Create a new Quick HTTP client with custom settings
	cClient := client.New(
		client.WithTimeout(5*time.Second),   // Sets the request timeout to 5 seconds
		client.WithDisableKeepAlives(false), // Enables persistent connections (Keep-Alive)
		client.WithMaxIdleConns(20),         // Defines a maximum of 20 idle connections
		client.WithMaxConnsPerHost(20),      // Limits simultaneous connections per host to 20
		client.WithMaxIdleConnsPerHost(20),  // Limits idle connections per host to 20
		client.WithContext(context.TODO()),  // Injects a context (can be used for cancellation)
		client.WithHeaders(
			map[string]string{
				"Content-Type":  "application/json", // Specifies the request content type
				"Authorization": "Bearer Token",     // Adds an authorization token for authentication
			},
		),
		client.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,             // ⚠ Disables SSL certificate verification (use with caution)
			MinVersion:         tls.VersionTLS12, // Enforces a minimum TLS version for security
		}),
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,                         // Allows up to 2 retry attempts for failed requests
				Delay:      1 * time.Second,           // Delay of 1 second between retries
				UseBackoff: true,                      // Enables exponential backoff for retries
				Statuses:   []int{502, 503, 504, 403}, // Retries only on specific HTTP status codes
				FailoverURLs: []string{ // Backup URLs in case the primary request fails
					"http://backup1",
					"https://reqres.in/api/users",
					"https://httpbin_error.org/post",
				},
				EnableLog: true, // Enables logging for debugging retry behavior
			}),
	)
	// Send a POST request to the primary URL
	resp, err := cClient.Post("https://httpbin_error.org/post",
		map[string]string{"message": "Hello, POST in Quick!"})
	if err != nil {
		log.Fatal(err) // Logs an error and exits if the request fails
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err) // Logs an error if the response cannot be parsed
	}

	// Print the response
	fmt.Println("POST response:", result)
}
