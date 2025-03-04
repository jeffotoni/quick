package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/jeffotoni/quick"
)

// ExampleClient_Get demonstrates using the Client's Get method.
func ExampleClient_Get() {
	// Create a test server that returns "GET OK" for GET requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET OK"))
	}))
	defer ts.Close()

	// Create a default client.
	c := New()

	// Send a GET request.
	resp, err := c.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// GET OK
}

// ExampleClient_Post demonstrates using the Client's Post method with a flexible body.
func ExampleClient_Post() {
	// Create a test server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}))
	defer ts.Close()

	// Create a default client.
	c := New()

	// Example 1: Using a string as the POST body.
	resp, err := c.Post(ts.URL, "Hello, POST!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Using a struct as the POST body (marshaled to JSON).
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, JSON POST!",
	}
	resp, err = c.Post(ts.URL, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println("Struct body:", result["message"])

	// Example 3: Using an io.Reader as the POST body.
	reader := strings.NewReader("Reader POST")
	resp, err = c.Post(ts.URL, reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("io.Reader body:", string(resp.Body))

	// Out put:
	// String body: Hello, POST!
	// Struct body: Hello, JSON POST!
	// io.Reader body: Reader POST
}

// ExampleClient_Put demonstrates using the Client's Put method with a flexible body.
func ExampleClient_Put() {
	// Create a test server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer ts.Close()

	// Create a default client.
	c := New()

	// Example 1: Using a string as the PUT body.
	resp, err := c.Put(ts.URL, "Hello, PUT!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Using a struct as the PUT body (marshaled to JSON).
	data := struct {
		Value int `json:"value"`
	}{Value: 42}
	resp, err = c.Put(ts.URL, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var result map[string]int
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println("Struct body:", result["value"])

	// Out put:
	// String body: Hello, PUT!
	// Struct body: 42
}

// ExampleClient_Delete demonstrates using the Client's Delete method.
func ExampleClient_Delete() {
	// Create a test server that returns "DELETE OK" for DELETE requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DELETE OK"))
	}))
	defer ts.Close()

	// Create a default client.
	c := New()

	// Send a DELETE request.
	resp, err := c.Delete(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// DELETE OK
}

// ExampleWithRetry demonstrates using the Client with retry logic.
func ExampleWithRetry() {
	// Create a test server that fails twice before succeeding.
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Retry succeeded"))
	}))
	defer ts.Close()

	// Create a client with retry enabled.
	c := New(WithRetry(RetryConfig{
		MaxRetries: 3,
		Delay:      500 * time.Millisecond,
		UseBackoff: true,
		Statuses:   []int{500},
		EnableLog:  false,
	}))

	resp, err := c.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Retry succeeded
}

// ExampleWithInsecureTLS demonstrates using the Client with insecure TLS enabled.
func ExampleWithInsecureTLS() {
	// Create a test TLS server.
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Insecure TLS OK"))
	}))
	defer ts.Close()

	// Create a client with insecure TLS enabled.
	c := New(WithInsecureTLS(true))
	resp, err := c.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Insecure TLS OK
}

// ExampleClient_full demonstrates full usage of the Client package with a server built using the quick framework,
// using extended client configuration options such as WithTimeout, WithDisableKeepAlives, WithMaxIdleConns,
// WithMaxConnsPerHost, WithMaxIdleConnsPerHost, WithContext, WithHeaders, WithCustomHTTPClient, WithTransportConfig,
// WithRetryRoundTripper, and WithRetry.
// Note: WithRetryRoundTripper overrides previous transport settings.
func ExampleClient_full() {
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
	client := New(
		WithTimeout(5*time.Second),
		WithDisableKeepAlives(false),
		WithMaxIdleConns(20),
		WithMaxConnsPerHost(20),
		WithMaxIdleConnsPerHost(20),
		WithContext(context.Background()),
		WithCustomHTTPClient(customHTTPClient),
		WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer EXAMPLE_TOKEN",
		}),
		// Set transport options via WithTransportConfig.
		WithTransportConfig(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     20,
			DisableKeepAlives:   false,
		}),
		// Override transport with retry RoundTripper.
		WithRetryRoundTripper(RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500},
			EnableLog:  false,
		}),
		// Also configure client retry settings (for manual retry logic).
		WithRetry(RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500},
			EnableLog:  false,
		}),
	)

	// GET request.
	resp, err := client.Get(ts.URL + "/get")
	if err != nil {
		fmt.Println("GET Error:", err)
		return
	}
	fmt.Println("GET:", string(resp.Body))

	// POST request with a string body.
	resp, err = client.Post(ts.URL+"/post", "Hello, extended POST!")
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
	resp, err = client.Put(ts.URL+"/put", data)
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
	resp, err = client.Delete(ts.URL + "/delete")
	if err != nil {
		fmt.Println("DELETE Error:", err)
		return
	}
	fmt.Println("DELETE:", string(resp.Body))

	// POSTFORM request.
	formData := url.Values{}
	formData.Set("key", "value")
	resp, err = client.PostForm(ts.URL+"/postform", formData)
	if err != nil {
		fmt.Println("POSTFORM Error:", err)
		return
	}
	fmt.Println("POSTFORM:", string(resp.Body))

	// Out put:
	// GET: GET OK
	// POST: POST: Hello, extended POST!
	// PUT: PUT: {"data":"Hello, extended PUT!"}
	// DELETE: DELETE OK
	// POSTFORM: POSTFORM: key=value
}
