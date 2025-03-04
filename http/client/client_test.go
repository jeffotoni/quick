package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// testHandler simulates responses for various HTTP methods.
func testHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET OK"))
	case http.MethodPost:
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DELETE OK"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// TestClient_Get verifies the GET method.
func TestClient_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	c := New() // use default client configuration
	resp, err := c.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != "GET OK" {
		t.Errorf("Expected body 'GET OK', got '%s'", string(resp.Body))
	}
}

// TestClient_Post verifies the POST method with various body types.
func TestClient_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	c := New()

	// Test with a string body.
	bodyStr := "Test POST"
	resp, err := c.Post(ts.URL, bodyStr)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if string(resp.Body) != bodyStr {
		t.Errorf("Expected body '%s', got '%s'", bodyStr, string(resp.Body))
	}

	// Test with a struct body (marshaled to JSON).
	type TestData struct {
		Message string `json:"message"`
	}
	data := TestData{Message: "Hello JSON"}
	resp, err = c.Post(ts.URL, data)
	if err != nil {
		t.Fatalf("POST request with struct failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	var result TestData
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if result.Message != data.Message {
		t.Errorf("Expected message '%s', got '%s'", data.Message, result.Message)
	}

	// Test with an io.Reader body.
	reader := strings.NewReader("Reader POST")
	resp, err = c.Post(ts.URL, reader)
	if err != nil {
		t.Fatalf("POST request with io.Reader failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if string(resp.Body) != "Reader POST" {
		t.Errorf("Expected body 'Reader POST', got '%s'", string(resp.Body))
	}
}

// TestClient_Put verifies the PUT method with various body types.
func TestClient_Put(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	c := New()

	// Test with a string body.
	bodyStr := "Test PUT"
	resp, err := c.Put(ts.URL, bodyStr)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != bodyStr {
		t.Errorf("Expected body '%s', got '%s'", bodyStr, string(resp.Body))
	}

	// Test with a struct body (marshaled to JSON).
	type TestData struct {
		Value int `json:"value"`
	}
	data := TestData{Value: 123}
	resp, err = c.Put(ts.URL, data)
	if err != nil {
		t.Fatalf("PUT request with struct failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var result TestData
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if result.Value != data.Value {
		t.Errorf("Expected value %d, got %d", data.Value, result.Value)
	}
}

// TestClient_Delete verifies the DELETE method.
func TestClient_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	c := New()
	resp, err := c.Delete(ts.URL)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != "DELETE OK" {
		t.Errorf("Expected body 'DELETE OK', got '%s'", string(resp.Body))
	}
}

// TestParseBody verifies the parseBody function with various input types.
func TestParseBody(t *testing.T) {
	// Test with nil input.
	r, err := parseBody(nil)
	if err != nil {
		t.Fatalf("parseBody(nil) failed: %v", err)
	}
	if r != nil {
		t.Errorf("Expected nil reader for nil body, got %v", r)
	}

	// Test with io.Reader input.
	original := "test"
	input := strings.NewReader(original)
	r, err = parseBody(input)
	if err != nil {
		t.Fatalf("parseBody(io.Reader) failed: %v", err)
	}
	buf, _ := io.ReadAll(r)
	if string(buf) != original {
		t.Errorf("Expected '%s', got '%s'", original, string(buf))
	}

	// Test with string input.
	r, err = parseBody("test string")
	if err != nil {
		t.Fatalf("parseBody(string) failed: %v", err)
	}
	buf, _ = io.ReadAll(r)
	if string(buf) != "test string" {
		t.Errorf("Expected 'test string', got '%s'", string(buf))
	}

	// Test with struct input (marshaled to JSON).
	type sample struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	s := sample{Name: "John", Age: 30}
	r, err = parseBody(s)
	if err != nil {
		t.Fatalf("parseBody(struct) failed: %v", err)
	}
	buf, _ = io.ReadAll(r)
	var result sample
	if err := json.Unmarshal(buf, &result); err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}
	if result != s {
		t.Errorf("Expected %+v, got %+v", s, result)
	}
}

// TestClient_WithCustomConfig verifies that a custom HTTPClientConfig is applied correctly.
func TestClient_WithCustomConfig(t *testing.T) {
	// Create a custom configuration with a shorter timeout and different transport settings.
	cfg := &HTTPClientConfig{
		Timeout:             5 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConns:        20,
		MaxConnsPerHost:     20,
		MaxIdleConnsPerHost: 20,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// Create a client with the custom configuration.
	c := New(
		WithContext(context.TODO()),
		WithHeaders(map[string]string{"Content-Type": "application/xml"}),
		WithHTTPClientConfig(cfg),
	)

	// Verify that the custom header is applied.
	if c.Headers["Content-Type"] != "application/xml" {
		t.Errorf("Expected header 'application/xml', got '%s'", c.Headers["Content-Type"])
	}

	// Verify that the custom timeout is applied.
	httpClient, ok := c.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatalf("ClientHTTP is not of type *http.Client")
	}
	if httpClient.Timeout != cfg.Timeout {
		t.Errorf("Expected timeout %v, got %v", cfg.Timeout, httpClient.Timeout)
	}
}

// TestClient_WithTimeout verifies that WithTimeout correctly sets the timeout and logging flag.
func TestClient_WithTimeout(t *testing.T) {
	c := New(
		WithTimeout(1 * time.Second),
	)
	httpClient, ok := c.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatalf("ClientHTTP is not of type *http.Client")
	}
	if httpClient.Timeout != 1*time.Second {
		t.Errorf("Expected timeout 1s, got %v", httpClient.Timeout)
	}
}
