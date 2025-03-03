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
	// Creating a CookieJar to manage cookies automatically
	jar, _ := cookiejar.New(nil)

	// Creating a fully custom *http.Client
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar, // Uses a CookieJar to store cookies
	}

	// Creating a quick client using the custom *http.Client
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Uses the pre-configured HTTP client
		client.WithContext(context.Background()),      // Sets a request context
		client.WithHeaders(map[string]string{
			"User-Agent": "QuickClient/1.0",
		}),
		client.WithRetry(3, "1s-bex", "500,502,503,504"), // Enables retry for specific HTTP status codes
	)

	resp, err := cClient.Post("http://localhost:3000/v1/user", map[string]string{"name": "jeffotoni"})
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	fmt.Println("POST Form Response:", string(resp.Body))
}
