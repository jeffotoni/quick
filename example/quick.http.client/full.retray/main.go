package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
// - WithTimeout: Sets the HTTP client timeout to 30 seconds.
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
func main() {

	// If both WithRetry RoundTripper and WithRetry are declared,
	// WithRetry RoundTripper takes precedence and is executed.
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Sets a global timeout for all requests
		client.WithHeaders(map[string]string{
			"Content-Type": "application/json", // Defines the content type for requests
		}),

		// Configures an advanced HTTP transport for connection optimization
		client.WithTransportConfig(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,             // Uses system proxy settings
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // Disables SSL verification (not recommended for production)
			ForceAttemptHTTP2:   true,                                  // Enables HTTP/2 if the server supports it
			MaxIdleConns:        20,                                    // Limits the total number of idle connections
			MaxIdleConnsPerHost: 10,                                    // Limits idle connections per host
			MaxConnsPerHost:     20,                                    // Limits concurrent connections per host
			DisableKeepAlives:   false,                                 // Enables persistent connections (Keep-Alive)
		}),

		// Enables automatic retries for failed requests
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,               // Maximum number of retry attempts
				Delay:      1 * time.Second, // Wait time before retrying
				UseBackoff: true,            // Enables exponential backoff for retries
				Statuses:   []int{500},      // Retries only on HTTP 500 errors
				FailoverURLs: []string{ // Backup URLs in case the primary request fails
					"http://backup1",
					"https://reqres.in/api/users",
					"https://httpbin_error.org/post",
				},
				EnableLog: true, // Enables logging for debugging retries
			}),
	)

	// Sends a POST request to the primary URL
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err) // Logs an error if the request fails
	}

	// Print the response body
	fmt.Println("POST Form Response:", string(resp.Body))
}
