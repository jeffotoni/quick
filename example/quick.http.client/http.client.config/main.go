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

	// Configuring the HTTP client using a structured approach.
	//
	// The following settings are applied to the HTTP client:
	// - Timeout: Sets the maximum duration for requests (20 seconds).
	// - DisableKeepAlives: Controls whether keep-alive connections are disabled (false = keep-alives enabled).
	// - MaxIdleConns: Defines the maximum number of idle connections across all hosts (20).
	// - MaxConnsPerHost: Sets the maximum number of simultaneous connections to a single host (20).
	// - MaxIdleConnsPerHost: Defines the maximum number of idle connections per host (20).
	// - TLSClientConfig: Configures TLS settings, including:
	//     * InsecureSkipVerify: false (enables strict TLS verification).
	//     * MinVersion: TLS 1.2 (ensures a minimum TLS version for security).
	//
	// Using WithHTTPClientConfig(cfg), all the configurations are applied at once.
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
