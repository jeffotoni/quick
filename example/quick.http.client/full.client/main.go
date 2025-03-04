package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Initialize the quick framework.
	q := quick.New()

	// Define routes using quick.
	q.Get("/get", func(c *quick.Ctx) error {
		return c.Status(200).SendString("GET OK")
	})
	q.Post("/post", func(c *quick.Ctx) error {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(201).SendString("POST: " + string(body))
	})
	q.Put("/put", func(c *quick.Ctx) error {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(200).SendString("PUT: " + string(body))
	})
	q.Delete("/delete", func(c *quick.Ctx) error {
		return c.Status(200).SendString("DELETE OK")
	})
	q.Post("/postform", func(c *quick.Ctx) error {
		// Assume FormValues returns map[string][]string.
		form := c.FormValues()
		vals := url.Values(form)
		return c.Status(200).SendString("POSTFORM: " + vals.Encode())
	})

	// Create a test server using the quick handler.
	ts := httptest.NewServer(q)
	defer ts.Close()

	// Creating a custom HTTP transport with advanced settings.
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // Uses system proxy settings if available.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,             // Allows insecure TLS connections (not recommended for production).
			MinVersion:         tls.VersionTLS12, // Enforces a minimum TLS version for security.
		},
		MaxIdleConns:        50,    // Maximum number of idle connections across all hosts.
		MaxConnsPerHost:     30,    // Maximum simultaneous connections per host.
		MaxIdleConnsPerHost: 10,    // Maximum number of idle connections per host.
		DisableKeepAlives:   false, // Enables persistent connections (Keep-Alive).
	}

	// Creating a fully custom *http.Client with the transport and timeout settings.
	customHTTPClient := &http.Client{
		Timeout:   15 * time.Second, // Global timeout for all requests.
		Transport: customTransport,  // Uses the custom transport.
	}

	// Create a client with extended options.
	cClient := client.New(
		client.WithTimeout(5*time.Second),
		client.WithDisableKeepAlives(false),
		client.WithMaxIdleConns(20),
		client.WithMaxConnsPerHost(20),
		client.WithMaxIdleConnsPerHost(20),
		client.WithContext(context.Background()),
		client.WithCustomHTTPClient(customHTTPClient),
		client.WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer EXAMPLE_TOKEN",
		}),
		// Set transport options via WithTransportConfig.

		// WithTransportConfig(&http.Transport{
		// 	Proxy:               http.ProxyFromEnvironment,
		// 	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		// 	ForceAttemptHTTP2:   true,
		// 	MaxIdleConns:        20,
		// 	MaxIdleConnsPerHost: 10,
		// 	MaxConnsPerHost:     20,
		// 	DisableKeepAlives:   false,
		// }),

		// Override transport with retry RoundTripper.
		// WithRetryRoundTripper(RetryConfig{
		// 	MaxRetries: 2,
		// 	Delay:      1 * time.Second,
		// 	UseBackoff: true,
		// 	Statuses:   []int{500},
		// 	EnableLog:  false,
		// }),

		// Also configure client retry settings (for manual retry logic).
		client.WithRetry(client.RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500},
			EnableLog:  false,
		}),
	)

	// GET request.
	resp, err := cClient.Get(ts.URL + "/get")
	if err != nil {
		fmt.Println("GET Error:", err)
		return
	}
	fmt.Println("GET:", string(resp.Body))

	// POST request with a string body.
	resp, err = cClient.Post(ts.URL+"/post", "Hello, extended POST!")
	if err != nil {
		fmt.Println("POST Error:", err)
		return
	}
	fmt.Println("POST:", string(resp.Body))

	// PUT request with a struct body (marshaled to JSON).
	data := struct {
		Data string `json:"data"`
	}{
		Data: "Hello, extended PUT!",
	}
	resp, err = cClient.Put(ts.URL+"/put", data)
	if err != nil {
		fmt.Println("PUT Error:", err)
		return
	}
	// To display the JSON response as a string, unmarshal and marshal it back.
	var putResult map[string]string
	_ = json.Unmarshal(resp.Body, &putResult)
	putJSON, _ := json.Marshal(putResult)
	fmt.Println("PUT:", string(putJSON))

	// DELETE request.
	resp, err = cClient.Delete(ts.URL + "/delete")
	if err != nil {
		fmt.Println("DELETE Error:", err)
		return
	}
	fmt.Println("DELETE:", string(resp.Body))

	// POSTFORM request.
	formData := url.Values{}
	formData.Set("key", "value")
	resp, err = cClient.PostForm(ts.URL+"/postform", formData)
	if err != nil {
		fmt.Println("POSTFORM Error:", err)
		return
	}
	fmt.Println("POSTFORM:", string(resp.Body))

	// Output:
	// GET: GET OK
	// POST: POST: Hello, extended POST!
	// PUT: PUT: {"data":"Hello, extended PUT!"}
	// DELETE: DELETE OK
	// POSTFORM: POSTFORM: key=value
}
