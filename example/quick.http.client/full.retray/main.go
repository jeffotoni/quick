package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// when I declare the 2 retrys, WithRetry RoundTripper and WithRetry ,
	// the With Retry RoundTripper overrides it which is executed.
	cClient := client.New(
		client.WithTimeout(10*time.Second),
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Advanced HTTP transport configuration
		client.WithRetryTransport(
			50,    // MaxIdleConns
			10,    // MaxIdleConnsPerHost
			30,    // MaxConnsPerHost
			false, // DisableKeepAlives
			true,  // Force HTTP/2
			http.ProxyFromEnvironment,
			&tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
			},
		),

		// Automatic retry via RoundTripper
		client.WithRetryRoundTripper(
			5,                    // Maximum number of retries
			1*time.Second,        // Delay between attempts
			true,                 // Use exponential backoff
			[]int{500, 502, 503}, // HTTP status for retry
		),

		// Retry quick
		// client.WithRetry(5, "2s-bex", "500,502,503,504"),
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Form Response:", string(resp.Body))
}
