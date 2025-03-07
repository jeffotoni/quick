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
// - WithContext: Injects a context for the client (context.TODO() used as placeholder).
// - WithHeaders: Adds custom headers (e.g., Content-Type: application/xml).
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
// - WithHTTPClientConfig: Defines advanced transport settings like connection pooling.
func main() {
	cfg := &client.HTTPClientConfig{
		Timeout:             20 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConns:        20,
		MaxConnsPerHost:     20,
		MaxIdleConnsPerHost: 20,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// Creating an HTTP client with the pre-defined configuration.
	//
	// - WithContext: Sets a custom context for handling request cancellation and deadlines.
	// - WithHeaders: Adds a map of default headers (e.g., "Content-Type: application/xml").
	// - WithHTTPClientConfig: Applies the entire configuration object (cfg) to the client.
	cClient := client.New(
		client.WithContext(context.TODO()),
		client.WithHeaders(map[string]string{"Content-Type": "application/xml"}),
		client.WithHTTPClientConfig(cfg),
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				FailoverURLs: []string{
					"http://backup1",
					"https://httpbin_error.org/post",
					"https://httpbin.org/post"},
				EnableLog: true,
			}),
	)

	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := cClient.Post("https://httpbin_error.org/post ", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response (if applicable)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result["message"])
}
