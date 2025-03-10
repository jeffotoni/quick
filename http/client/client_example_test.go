package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

// ExampleClient_Get demonstrates how to use the Client's Get method.
func ExampleClient_Get() {
	// Create a test HTTP server that responds with "GET OK" to GET requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET OK"))
	}))
	defer ts.Close()

	// Initialize a new client.
	client := New()

	// Perform a GET request to the test server.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// GET OK
}

// ExampleClient_Post demonstrates how to use the Client's Post method with different types of request bodies.
func ExampleClient_Post() {
	// Create a test HTTP server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}))
	defer ts.Close()

	// Initialize a new client.
	client := New()

	// Example 1: Sending a string as the POST body.
	resp, err := client.Post(ts.URL, "Hello, POST!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Sending a struct as the POST body (automatically marshaled to JSON).
	data := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, JSON POST!",
	}
	resp, err = client.Post(ts.URL, data)
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

	// Example 3: Sending an io.Reader as the POST body.
	reader := strings.NewReader("Reader POST")
	resp, err = client.Post(ts.URL, reader)
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

// ExampleClient_Put demonstrates how to use the Client's Put method with different types of request bodies.
func ExampleClient_Put() {
	// Create a test HTTP server that echoes the request body.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer ts.Close()

	// Initialize a new client.
	client := New()

	// Example 1: Sending a string as the PUT body.
	resp, err := client.Put(ts.URL, "Hello, PUT!")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("String body:", string(resp.Body))

	// Example 2: Sending a struct as the PUT body (automatically marshaled to JSON).
	data := struct {
		Value int `json:"value"`
	}{Value: 42}

	resp, err = client.Put(ts.URL, data)
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

// ExampleClient_Delete demonstrates how to use the Client's Delete method.
func ExampleClient_Delete() {
	// Create a test HTTP server that responds with "DELETE OK" to DELETE requests.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DELETE OK"))
	}))
	defer ts.Close()

	// Initialize a new client.
	client := New()

	// Perform a DELETE request to the test server.
	resp, err := client.Delete(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// DELETE OK
}

// ExampleWithInsecureTLS demonstrates how to use the Client with insecure TLS enabled.
func ExampleWithInsecureTLS() {
	// Create a test TLS server that returns a response for HTTPS requests.
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Insecure TLS OK"))
	}))
	defer ts.Close()

	// Initialize a client with insecure TLS verification enabled.
	client := New(WithInsecureTLS(true))

	// Perform a GET request to the TLS test server.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Insecure TLS OK
}

// ExampleWithCustomHTTPClient demonstrates how to use the Client with a fully custom *http.Client.
func ExampleWithCustomHTTPClient() {
	// Create a test HTTP server that responds with "Custom Client OK".
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Custom Client OK"))
	}))
	defer ts.Close()

	// Create a custom HTTP client with specific settings.
	customHTTPClient := &http.Client{
		Timeout: 5 * time.Second, // Set a custom timeout.
		Transport: &http.Transport{
			MaxIdleConns:        50, // Limit the maximum idle connections.
			MaxConnsPerHost:     20, // Limit concurrent connections per host.
			MaxIdleConnsPerHost: 10, // Limit idle connections per host.
		},
	}

	// Initialize a Client with the custom HTTP client.
	client := New(WithCustomHTTPClient(customHTTPClient))

	// Perform a GET request using the custom client.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Output:
	// Custom Client OK
}

// ExampleWithTimeout demonstrates how to use the Client with a custom timeout.
func ExampleWithTimeout() {
	// Create a test HTTP server that delays its response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delayed response"))
	}))
	defer ts.Close()

	// Initialize a client with a 50ms timeout.
	client := New(WithTimeout(50 * time.Millisecond))

	// Perform a GET request (expected to timeout).
	_, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Request timed out")
	} else {
		fmt.Println("Request succeeded")
	}

	// Out put:
	// Request timed out
}

// ExampleWithContext demonstrates how to use the Client with a custom context.
func ExampleWithContext() {
	// Create a test HTTP server that delays its response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delayed response"))
	}))
	defer ts.Close()

	// Create a context with a 100ms deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Initialize a client with the custom context.
	client := New(WithContext(ctx))

	// Perform a GET request (expected to timeout).
	_, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Request canceled due to timeout")
	} else {
		fmt.Println("Request succeeded")
	}

	// Out put:
	// Request canceled due to timeout
}

// ExampleWithHeaders demonstrates how to use the Client with custom headers.
func ExampleWithHeaders() {
	// Create a test HTTP server that checks for a custom header.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") == "GoLang" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Header received"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	// Initialize a client with a custom header.
	client := New(WithHeaders(map[string]string{"X-Custom-Header": "GoLang"}))

	// Perform a GET request.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Header received
}

// ExampleWithLogger demonstrates how to use the Client with logging enabled.
func ExampleWithLogger() {
	// Create a log buffer to capture logs.
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Initialize a client with logging enabled.
	client := New(WithLogger(true))
	client.Logger = logger

	// Perform a log entry.
	client.log("Testing logging")

	// Print captured log.
	fmt.Println(logBuffer.String())

	// Out put:
	// time=<timestamp> level=INFO msg="Testing logging"
}

// ExampleWithRetry demonstrates how to use the Client with retry logic.
func ExampleWithRetry() {
	// Create a test HTTP server that fails twice before succeeding.
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

	// Initialize a client with retry logic enabled.
	client := New(WithRetry(RetryConfig{
		MaxRetries: 3,
		Delay:      500 * time.Millisecond,
		UseBackoff: true,
		Statuses:   []int{http.StatusInternalServerError},
		EnableLog:  false,
	}))

	// Perform a GET request to the test server.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Retry succeeded
}

// ExampleNew_withRetry demonstrates how to use the Client with retry logic.
func ExampleNew_withRetry() {
	// Create a test HTTP server that fails twice before succeeding.
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

	// Initialize a client with retry logic enabled.
	client := New(
		WithRetry(RetryConfig{
			MaxRetries:   3,                                                                   // Allow up to 3 retry attempts.
			Delay:        500 * time.Millisecond,                                              // Base delay before retrying.
			UseBackoff:   true,                                                                // Use exponential backoff for retries.
			Statuses:     []int{http.StatusInternalServerError},                               // Retry on HTTP 500 errors.
			FailoverURLs: []string{"https://reqres.in/api/users", "https://httpbin.org/post"}, // Failover URLs.
			EnableLog:    false,                                                               // Disable logging for this example.
		}))

	// Perform a GET request to the test server.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Retry succeeded
}

// ExampleWithTLSConfig demonstrates how to use the Client with a custom TLS configuration.
func ExampleWithTLSConfig() {
	// Create a test TLS server.
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Secure TLS OK"))
	}))
	defer ts.Close()

	// Custom TLS configuration that allows insecure certificates (for testing purposes).
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// Initialize a client with the custom TLS configuration.
	client := New(WithTLSConfig(tlsConfig))

	// Perform a GET request to the TLS test server.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Secure TLS OK
}

// ExampleWithTransport demonstrates how to use the Client with a custom transport.
func ExampleWithTransport() {
	// Create a test HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Custom Transport OK"))
	}))
	defer ts.Close()

	// Create a custom transport.
	customTransport := &http.Transport{MaxIdleConns: 100}

	// Initialize a client with the custom transport.
	client := New(WithTransport(customTransport))

	// Perform a GET request.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Custom Transport OK
}

// ExampleWithTransportConfig demonstrates how to use the Client with a custom transport configuration.
func ExampleWithTransportConfig() {
	// Create a test HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Configured Transport OK"))
	}))
	defer ts.Close()

	// Pre-configured transport.
	preConfiguredTransport := &http.Transport{MaxConnsPerHost: 50}

	// Initialize a client with the pre-configured transport.
	client := New(WithTransportConfig(preConfiguredTransport))

	// Perform a GET request.
	resp, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(resp.Body))

	// Out put:
	// Configured Transport OK
}

// ExampleClient_PostForm demonstrates how to use the Client's PostForm method to send URL-encoded form data.
// This function is named ExampleClient_PostForm()
// it with the Examples type.
func ExampleClient_PostForm() {
	// Creating a test server that returns the form data.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm() // Analyzes the form data
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Form.Encode())) // Returns the form data already formatted
	}))
	defer ts.Close()

	// Initializing the client.
	client := New()

	// Creating the form data.
	formData := url.Values{
		"key":   {"value"},
		"hello": {"world"},
	}

	// Enviando a requisição POST com o formulário.
	resp, err := client.PostForm(ts.URL, formData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Displaying the answer.
	fmt.Println("Form POST response:", string(resp.Body))

	// Out put:
	// Form POST response: hello=world&key=value
}

// ExampleWithMaxIdleConns demonstrates how to configure the maximum number of idle connections in the HTTP client.
// This function is named ExampleWithMaxIdleConns()
// it with the Examples type.
func ExampleWithMaxIdleConns() {
	// Create a test HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()

	// Initialize a client with a custom maximum number of idle connections.
	client := New(WithMaxIdleConns(20))

	// Perform a GET request.
	_, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// You can print something related to the configuration to avoid the no output issue, or simply avoid the output comment
	fmt.Println("MaxIdleConns set to 20")

	// This will now match with the Output comment below
	// Out put:
	// MaxIdleConns set to 20
}

// ExampleWithMaxConnsPerHost demonstrates how to set the maximum number of concurrent connections per host for the HTTP client.
// This function is named ExampleWithMaxConnsPerHost()
// it with the Examples type.
func ExampleWithMaxConnsPerHost() {
	// Create a test HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()

	// Initialize a client with a custom maximum number of connections per host.
	client := New(WithMaxConnsPerHost(10))

	// Perform a GET request.
	_, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Confirm the setting (for example purposes, normally you might not print this).
	fmt.Println("MaxConnsPerHost set to 10")

	// Out put:
	// MaxConnsPerHost set to 10
}

// ExampleWithDisableKeepAlives demonstrates how to disable HTTP keep-alive connections.
// This function is named ExampleWithMaxConnsPerHost()
// it with the Examples type.
func ExampleWithDisableKeepAlives() {
	// Create a test HTTP server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()

	// Initialize a client with keep-alives disabled.
	client := New(WithDisableKeepAlives(true))

	// Perform a GET request.
	_, err := client.Get(ts.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Confirm that keep-alives are disabled (for example purposes, normally you might not print this).
	fmt.Println("Keep-alives are disabled")

	// Out put:
	// Keep-alives are disabled
}

// ExampleRetryTransport_RoundTrip demonstrates how RetryTransport handles retries and failover logic.
// This function is named ExampleRetryTransport_RoundTrip()
// it with the Examples type.
func ExampleRetryTransport_RoundTrip() {
	// Simulate an HTTP server that fails twice before responding successfully.
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Returns HTTP error 500 on the first two attempts.
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Temporary Server Error"))
			return
		}
		// On the third attempt, returns success.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success on retry!"))
	}))
	defer ts.Close()

	// Create a transport with retry and failover policy.
	retryTransport := &RetryTransport{
		Base:          http.DefaultTransport,
		MaxRetries:    3,                                     // Attempts up to 3 times before failing.
		RetryDelay:    200 * time.Millisecond,                // Wait time between attempts.
		UseBackoff:    true,                                  // Enables exponential backoff.
		RetryStatuses: []int{http.StatusInternalServerError}, // Retries only on HTTP 500 status.
		FailoverURLs:  []string{},                            // No alternative URLs.
		EnableLogger:  true,
		Logger:        slog.Default(),
	}

	// Create an HTTP client using our RetryTransport.
	client := &http.Client{Transport: retryTransport}

	// Create a test request.
	req, _ := http.NewRequest("GET", ts.URL, nil)

	// Execute the request.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response.
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))

	// Output:
	// Response: Success on retry!
}
