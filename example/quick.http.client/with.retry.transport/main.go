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
//   - WithTimeout: Sets the HTTP client timeout to 30 seconds.
//   - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
//   - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
//   - WithTransportConfig: Configures transport settings like max connections and keep-alives.
func main() {
	// Create an HTTP client with custom configurations using the Quick framework.
	cClient := client.New(
		// Set a global timeout for all requests made by this client to 10 seconds.
		// This helps prevent the client from hanging indefinitely on requests.
		client.WithTimeout(10*time.Second),

		// Set default headers for all requests made by this client.
		// Here, we specify that we expect to send and receive JSON data.
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Configure the underlying transport for the HTTP client.
		client.WithTransportConfig(&http.Transport{
			// Use the system environment settings for proxy configuration.
			Proxy: http.ProxyFromEnvironment,

			// Configure TLS settings to skip verification of the server's certificate chain and hostname.
			// Warning: Setting InsecureSkipVerify to true is not recommended for production as it is insecure.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},

			// Enable HTTP/2 for supported servers.
			ForceAttemptHTTP2: true,

			// Set the maximum number of idle connections in the connection pool for all hosts.
			MaxIdleConns: 20,

			// Set the maximum number of idle connections in the connection pool per host.
			MaxIdleConnsPerHost: 10,

			// Set the maximum number of simultaneous connections per host.
			MaxConnsPerHost: 20,

			// Keep connections alive between requests. This can help improve performance.
			DisableKeepAlives: false,
		}),
	)

	// Perform a POST request with a JSON payload.
	// The payload includes a single field "name" with a value.
	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		// Log the error and stop the program if the POST request fails.
		log.Fatalf("POST request failed: %v", err)
	}

	// Output the response from the POST request.
	fmt.Println("POST Form Response:", string(resp.Body))
}
