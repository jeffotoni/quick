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
		client.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
		// Enabling retry
		client.WithRetry(
			client.RetryConfig{
				MaxRetries:   2,
				Delay:        1 * time.Second,
				UseBackoff:   true,
				Statuses:     []int{500, 502, 503},
				FailoverURLs: []string{"http://hosterror", "https://httpbin.org/post"},
				EnableLog:    true,
			}),
	)

	// call POST
	resp, err := cClient.Post("https://httpbin_error.org/post", map[string]string{
		"name": "jeffotoni in action with Quick!!!",
	})

	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Form Response:", string(resp.Body))
}
