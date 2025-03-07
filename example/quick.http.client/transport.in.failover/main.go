package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// Creating a custom HTTP transport with advanced settings.
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // Uses system proxy settings if available.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,             // Allows insecure TLS connections (not recommended for production).
			MinVersion:         tls.VersionTLS12, // Enforces a minimum TLS version for security.
		},
		MaxIdleConns:        50,    // Maximum number of idle connections across all hosts.
		MaxConnsPerHost:     30,    // Maximum simultaneous connections per host.
		MaxIdleConnsPerHost: 10,    // Maximum number of idle connections per host.
		DisableKeepAlives:   false, // Enables persistent connections (Keep-Alive).
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		Timeout:   5 * time.Second, // Sets a global timeout for all requests.
		Transport: customTransport, // Uses the custom transport.
	}

	// Creating a client using both the custom transport and other configurations.
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Applying the custom HTTP client.
		client.WithContext(context.Background()),      // Custom context for request cancellation and deadlines.
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		client.WithTimeout(5*time.Second), // Setting a timeout for requests.
		// Retry on specific status codes.
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500},
				FailoverURLs: []string{"https://httpbin_error.org/post", "https://httpbin.org/post"},
				EnableLog:    true,
			}),
	)

	resp, err := cClient.Post("https://httpbin.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	// show resp
	fmt.Println("POST response:", string(resp.Body))
}
