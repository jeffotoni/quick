package client

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// TestWithRetryTransport verifies that the transport settings are correctly applied.
func TestWithRetryTransport(t *testing.T) {
	cClient := New(
		WithRetryTransport(
			50,    // MaxIdleConns
			10,    // MaxIdleConnsPerHost
			30,    // MaxConnsPerHost
			false, // DisableKeepAlives
			true,  // Force HTTP/2
			http.ProxyFromEnvironment,
			&tls.Config{InsecureSkipVerify: true},
		),
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

	slogerdefult := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cClient := New(
		WithRetryRoundTripper(3, 2*time.Second, true, []int{500, 502, 503}, slogerdefult),
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
		w.WriteHeader(http.StatusServiceUnavailable) // Always return 503
	}))
	defer server.Close()

	slogerdefult := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	retryTransport := &RetryTransport{
		Base:        http.DefaultTransport,
		MaxRetries:  2,
		RetryDelay:  100 * time.Millisecond,
		UseBackoff:  true,
		RetryStatus: []int{503},
		Logger:      slogerdefult, // New Logger field
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
