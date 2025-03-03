package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

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
	// - "2mil-bex"  => 2 milliseconds with exponential backoff
	// - "2s-bex"    => 2 seconds with exponential backoff
	// - "2min-bex"  => 2 minutes with exponential backoff
	//
	// When using the "-bex" suffix, the delay will increase exponentially on each retry attempt
	// using the formula: waitTime = baseDelay * (2^attempt).
	//
	// Example Usage:
	//
	// client.WithRetry(3, "2s-bex", "500,502,503,504")
	//
	// This configuration will retry up to 3 times with an exponential backoff starting at 2 seconds,
	// and will only retry if the response status is 500, 502, 503, or 504.
	cClient := client.New(
		client.WithRetry(3, "2s-bex", "500,502,503,504"),
	)

	// Define the data payload for the POST request
	data := map[string]string{
		"name":  "Jefferson",
		"email": "jeff@example.com",
	}

	// Perform the POST request
	resp, err := cClient.Post("http://localhost:3000/v1/user", data)
	if err != nil {
		log.Fatal("Error making POST request:", err)
	}

	// Print the response body and status code
	fmt.Println("POST Response Status:", resp.StatusCode)
	fmt.Println("POST Response Body:", string(resp.Body))
}
