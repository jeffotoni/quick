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
			MaxRetries: 3,                         // Maximum number of retries.
			Delay:      2 * time.Second,           // Delay between attempts.
			UseBackoff: true,                      // Use exponential backoff.
			Statuses:   []int{500, 502, 503, 504}, // HTTP statuses for retry.
			EnableLog:  true,                      // Enable logging.
		}),
	)

	// Making a POST request.
	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	fmt.Println("POST Form Response:", string(resp.Body))
}
