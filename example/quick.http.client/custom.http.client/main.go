package main

import (
	"context"
	"crypto/tls"
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
		Timeout: 10 * time.Second, // Sets a global timeout of 10 seconds
		Jar:     jar,              // Uses a CookieJar to store cookies
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allows up to 3 redirects
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Allows insecure TLS (not recommended for production)
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        50,    // Maximum idle connections
			MaxConnsPerHost:     30,    // Max simultaneous connections per host
			MaxIdleConnsPerHost: 10,    // Max idle connections per host
			DisableKeepAlives:   false, // Enables keep-alive
		},
	}

	// Creating a quick client using the custom *http.Client
	cClient := client.New(
		client.WithCustomHTTPClient(customHTTPClient), // Uses the pre-configured HTTP client
		client.WithContext(context.Background()),      // Sets a request context
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_ACCESS_TOKEN",
		}),
		client.WithRetry(3, "1s-bex", "500,502,503,504"), // Enables retry for specific HTTP status codes
	)

	// Performing a GET request
	resp, err := cClient.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	fmt.Println("GET Response:", string(resp.Body))

	// Performing a POST request
	data := map[string]string{"name": "QuickFramework", "version": "1.0"}
	resp, err = cClient.Post("https://httpbin.org/post", data)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	fmt.Println("POST Response:", string(resp.Body))
}
