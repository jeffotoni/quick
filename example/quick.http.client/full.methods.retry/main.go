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

func main() {

	// Example of creating an HTTP client using a fluent and modular approach.
	// This allows fine-grained control over HTTP settings without requiring a full config struct.
	//
	// - WithTimeout: Sets the HTTP client timeout to 30 seconds.
	// - WithDisableKeepAlives: Enables or disables HTTP keep-alives (false = keep-alives enabled).
	// - WithMaxIdleConns: Defines the maximum number of idle connections (20).
	// - WithMaxConnsPerHost: Sets the maximum connections allowed per host (20).
	// - WithMaxIdleConnsPerHost: Sets the maximum number of idle connections per host (20).
	// - WithContext: Injects a context for the client (context.TODO() used as placeholder).
	// - WithHeaders: Adds custom headers (e.g., Content-Type: application/xml).
	// - WithTLSConfig: Configures TLS settings, including InsecureSkipVerify and TLS version.
	// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
	//   with exponential backoff (2s-bex) and a maximum of 3 attempts.
	cClient := client.New(
		client.WithTimeout(30*time.Second),
		client.WithDisableKeepAlives(false),
		client.WithMaxIdleConns(20),
		client.WithMaxConnsPerHost(20),
		client.WithMaxIdleConnsPerHost(20),
		client.WithContext(context.TODO()),
		client.WithHeaders(map[string]string{"Content-Type": "application/xml"}),
		client.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}),
		client.WithRetry(3, "2s-bex", "500,502,503,504"),
	)

	// Define a struct to send as JSON
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, POST!",
	}

	resp, err := cClient.Post("http://localhost:3000/v1/user", data)
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
