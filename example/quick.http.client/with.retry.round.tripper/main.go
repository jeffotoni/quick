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
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				EnableLog:  true,
			}),
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Form Response:", string(resp.Body))
}
