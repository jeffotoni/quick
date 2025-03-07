package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

// Example of creating an HTTP client using a fluent and modular approach.
// This allows fine-grained control over HTTP settings without requiring a full config struct.
//
// - WithRetry: Enables automatic retries for specific HTTP status codes (500, 502, 503, 504)
func main() {

	// Retry Delay Format Support
	//
	// The retry delay parameter supports various formats for flexibility in defining the wait time between retries.
	// Additionally, it allows enabling exponential backoff by appending "-bex".
	//
	// Supported formats:
	//
	// - "2mil"      => 2 milliseconds
	// - "2s"        => 2 seconds
	// - "2min"      => 2 minutes
	//
	// Example Usage:
	//
	// client.WithRetry(
	// 		client.RetryConfig{
	// 			MaxRetries: 2,
	// 			Delay:      1 * time.Second,
	// 			UseBackoff: true,
	// 			Statuses:   []int{500},
	// 			EnableLog:  true,
	// 		}),
	//
	// This configuration will retry up to 3 times with an exponential backoff starting at 2 seconds,
	// and will only retry if the response status is 500, 502, 503, or 504.
	cClient := client.New(
		client.WithRetry(
			client.RetryConfig{
				MaxRetries: 2,
				Delay:      1 * time.Second,
				UseBackoff: true,
				Statuses:   []int{500},
				FailoverURLs: []string{
					"http://backup1",
					"https://reqres.in/api/users",
					"https://httpbin.org/get"},
				EnableLog: true,
			}),
	)

	resp, err := cClient.Get("https://httpbin_error.org/get")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}
