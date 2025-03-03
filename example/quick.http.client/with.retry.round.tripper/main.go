package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := client.New(
		client.WithTimeout(10*time.Second),
		client.WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Enabling retry with RoundTripper
		client.WithRetryRoundTripper(
			3,                    // Maximum number of retries
			2*time.Second,        // Delay between retries
			true,                 // Use exponential backoff
			[]int{500, 502, 503}, // Retry these HTTP codes
		),
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Form Response:", string(resp.Body))
}
