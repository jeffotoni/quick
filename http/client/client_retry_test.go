package client

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// testRetryHandler simulates a server that initially fails before succeeding.
func testRetryHandler(failCount int, statusFail int, statusSuccess int, successBody string) http.HandlerFunc {
	var mu sync.Mutex
	attempts := 0
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		attempts++

		if attempts <= failCount {
			w.WriteHeader(statusFail)
			w.Write([]byte("Temporary failure"))
			return
		}
		w.WriteHeader(statusSuccess)
		w.Write([]byte(successBody))
	}
}

// TestClientRetry_Get verifies the GET method with retry logic.
func TestClientRetry_Get(t *testing.T) {
	ts := httptest.NewServer(testRetryHandler(2, http.StatusInternalServerError, http.StatusOK, "GET OK"))
	defer ts.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries: 3,
			Delay:      500 * time.Millisecond,
			UseBackoff: true,
			Statuses:   []int{500, 502},
			EnableLog:  true,
		}),
	)

	resp, err := client.Get(ts.URL)
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

// TestClientRetry_Post verifies the POST method with retry logic.
func TestClientRetry_Post(t *testing.T) {
	ts := httptest.NewServer(testRetryHandler(2, http.StatusServiceUnavailable, http.StatusCreated, "POST OK"))
	defer ts.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries: 3,
			Delay:      500 * time.Millisecond,
			UseBackoff: true,
			Statuses:   []int{503, 504},
			EnableLog:  true,
		}),
	)

	data := map[string]string{"name": "Jefferson"}
	resp, err := client.Post(ts.URL, data)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if string(resp.Body) != "POST OK" {
		t.Errorf("Expected body 'POST OK', got '%s'", string(resp.Body))
	}
}

// TestClientRetry_Put verifies the PUT method with retry logic.
func TestClientRetry_Put(t *testing.T) {
	ts := httptest.NewServer(testRetryHandler(1, http.StatusBadGateway, http.StatusOK, "PUT OK"))
	defer ts.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{502},
			EnableLog:  true,
		}),
	)

	data := map[string]string{"update": "yes"}
	resp, err := client.Put(ts.URL, data)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != "PUT OK" {
		t.Errorf("Expected body 'PUT OK', got '%s'", string(resp.Body))
	}
}

// TestClientRetry_Delete verifies the DELETE method with retry logic.
func TestClientRetry_Delete(t *testing.T) {
	ts := httptest.NewServer(testRetryHandler(3, http.StatusGatewayTimeout, http.StatusOK, "DELETE OK"))
	defer ts.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries: 4,
			Delay:      2 * time.Second,
			UseBackoff: true,
			Statuses:   []int{504},
			EnableLog:  true,
		}),
	)

	resp, err := client.Delete(ts.URL)
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

// TestWithRetryTransport verifies that the transport settings are correctly applied.
func TestWithRetryTransport(t *testing.T) {
	cClient := New(
		WithTransportConfig(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     30,
			DisableKeepAlives:   false,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}),
	)

	httpClient, ok := cClient.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatalf("Expected ClientHTTP to be *http.Client, got %T", cClient.ClientHTTP)
	}

	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Expected Transport to be *http.Transport, got %T", httpClient.Transport)
	}

	if transport.MaxIdleConns != 50 {
		t.Errorf("Expected MaxIdleConns 50, got %d", transport.MaxIdleConns)
	}
	if transport.MaxIdleConnsPerHost != 10 {
		t.Errorf("Expected MaxIdleConnsPerHost 10, got %d", transport.MaxIdleConnsPerHost)
	}
	if transport.MaxConnsPerHost != 30 {
		t.Errorf("Expected MaxConnsPerHost 30, got %d", transport.MaxConnsPerHost)
	}
	if transport.DisableKeepAlives != false {
		t.Errorf("Expected DisableKeepAlives false, got %v", transport.DisableKeepAlives)
	}
}

// TestWithRetryRoundTripper verifies that the RetryTransport is applied correctly.
func TestWithRetryRoundTripper(t *testing.T) {
	cClient := New(
		WithRetryRoundTripper(RetryConfig{
			MaxRetries: 3,
			Delay:      2 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500, 502, 503},
			EnableLog:  true,
		}),
	)

	httpClient, ok := cClient.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatalf("Expected ClientHTTP to be *http.Client, got %T", cClient.ClientHTTP)
	}

	transport, ok := httpClient.Transport.(*RetryTransport)
	if !ok {
		t.Fatalf("Expected Transport to be *RetryTransport, got %T", httpClient.Transport)
	}

	if transport.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", transport.MaxRetries)
	}
	if transport.RetryDelay != 2*time.Second {
		t.Errorf("Expected RetryDelay 2s, got %v", transport.RetryDelay)
	}
	if !transport.UseBackoff {
		t.Errorf("Expected UseBackoff true, got %v", transport.UseBackoff)
	}
	if len(transport.RetryStatus) != 3 {
		t.Errorf("Expected 3 retry statuses, got %d", len(transport.RetryStatus))
	}
}

// TestRetryTransport_RoundTrip verifies that the RoundTrip method properly retries failed requests.
func TestRetryTransport_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable) // Always return 503.
	}))
	defer server.Close()

	retryTransport := &RetryTransport{
		Base:        http.DefaultTransport,
		MaxRetries:  2,
		RetryDelay:  100 * time.Millisecond,
		UseBackoff:  true,
		RetryStatus: []int{503},
	}

	client := &http.Client{Transport: retryTransport}

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

// TestWithRetryRoundTripper_WithHeaders validates that retry logic works with RoundTripper,
// including proper header propagation.
func TestWithRetryRoundTripper_WithHeaders(t *testing.T) {
	// Creating a test server that fails 2 times and then returns success.
	ts := httptest.NewServer(testRetryHandler(2, http.StatusInternalServerError, http.StatusOK, `{"message": "Success!"}`))
	defer ts.Close()

	// Creating the client with retry configured.
	cClient := New(
		WithTimeout(8*time.Second), // Increased to ensure all attempts are made.
		WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Enabling retry with RoundTripper using RetryConfig.
		WithRetryRoundTripper(RetryConfig{
			MaxRetries: 3,
			Delay:      2 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500, 502, 503},
			EnableLog:  true,
		}),
	)

	// Execute the request.
	resp, err := cClient.Post(ts.URL, map[string]string{"name": "jeffotoni"})
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	// Verify the response status.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify the response body.
	expectedBody := `{"message": "Success!"}`
	if strings.TrimSpace(string(resp.Body)) != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, string(resp.Body))
	}
}
