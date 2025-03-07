package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
// - WithTimeout: Sets the HTTP client timeout to 30 seconds.
// - WithContext: Injects a context for the client (context.TODO() used as placeholder).
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/json).
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
// - WithCustomHTTPClient: Uses a pre-configured http.Client with custom settings.
func main() {
	// Creating a CookieJar to manage cookies automatically.
	jar, _ := cookiejar.New(nil)

	// Creating a fully custom *http.Client.
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second, // Sets a global timeout of 10 seconds.
		Jar:     jar,              // Uses a CookieJar to store cookies.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allows up to 3 redirects.
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Allows insecure TLS (not recommended for production).
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        50,    // Maximum idle connections.
			MaxConnsPerHost:     30,    // Max simultaneous connections per host.
			MaxIdleConnsPerHost: 10,    // Max idle connections per host.
			DisableKeepAlives:   false, // Enables keep-alive.
		},
	}

	// Creating a quick client using the custom *http.Client.
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Uses the pre-configured HTTP client.
		client.WithContext(context.Background()),      // Sets a request context.
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		// Enables retry for specific HTTP status codes using the new RetryConfig.
		client.WithRetry(client.RetryConfig{
			MaxRetries:   3,                         // Maximum number of retries.
			Delay:        1 * time.Second,           // Delay between attempts.
			UseBackoff:   true,                      // Use exponential backoff.
			Statuses:     []int{500, 502, 503, 504}, // HTTP statuses for retry.
			FailoverURLs: []string{"https://httpbin_error.org/get", "https://httpbin.org/get"},
			EnableLog:    true, // Enable logger.
		}),
	)

	// Creating a quick client using the custom *http.Client.
	cClient2 := client.New(
		client.WithTimeout(5*time.Second),
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		// Enables retry for specific HTTP status codes using the new RetryConfig.
		client.WithRetry(client.RetryConfig{
			MaxRetries:   3,                         // Maximum number of retries.
			Delay:        1 * time.Second,           // Delay between attempts.
			UseBackoff:   true,                      // Use exponential backoff.
			Statuses:     []int{500, 502, 503, 504}, // HTTP statuses for retry.
			FailoverURLs: []string{"https://httpbin_error.org/post", "https://httpbin.org/post"},
			EnableLog:    true, // Enable logger.
		}),
	)

	// Performing a GET request.
	resp, err := cClient.Get("https://httpbin_error.org/get")
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	fmt.Println("GET Response:", string(resp.Body))

	// Performing a POST request.
	data := map[string]string{"name": "QuickFramework", "version": "1.0"}
	resp, err = cClient2.Post("https://httpbin_error.org/post", data)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Response:", string(resp.Body))
}
