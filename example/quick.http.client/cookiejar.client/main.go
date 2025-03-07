package main

import (
	"context"
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
// - WithContext: Injects a context for the client (context.TODO() used as placeholder).
// - WithHeaders: Adds custom headers (e.g., User-Agent: QuickClient/1.0).
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
// - WithCustomHTTPClient: Uses a pre-configured http.Client with custom settings.
func main() {
	// Creating a CookieJar to manage cookies automatically.
	jar, _ := cookiejar.New(nil)

	// Creating a fully custom *http.Client.
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar, // Uses a CookieJar to store cookies.
	}

	// Creating a quick client using the custom *http.Client.
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Uses the pre-configured HTTP client.
		client.WithContext(context.Background()),      // Sets a request context.
		client.WithHeaders(map[string]string{
			"User-Agent": "QuickClient/1.0",
		}),
		client.WithRetry(client.RetryConfig{
			MaxRetries:   3,                         // Maximum number of retries.
			Delay:        2 * time.Second,           // Delay between attempts.
			UseBackoff:   true,                      // Use exponential backoff.
			Statuses:     []int{500, 502, 503, 504}, // HTTP statuses for retry.
			FailoverURLs: []string{"http://hosterror", "https://httpbin.org/post"},
			EnableLog:    true, // Enable logging.
		}),
	)

	// Making a POST request.
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	fmt.Println("POST Form Response:", string(resp.Body))
}
