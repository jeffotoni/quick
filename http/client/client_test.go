package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer ts.Close()

	client := New()
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(string(resp.Body), "success") {
		t.Errorf("Expected response body to contain 'success', got: %s", string(resp.Body))
	}
}

func TestRetryLogic(t *testing.T) {
	var attempt int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries: 3,
			Delay:      10 * time.Millisecond,
			UseBackoff: false,
			Statuses:   []int{http.StatusServiceUnavailable},
			EnableLog:  false,
		}),
	)

	_, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if attempt != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempt)
	}
}

func TestFailover(t *testing.T) {
	var primaryCalled, secondaryCalled bool

	primary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		primaryCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer primary.Close()

	secondary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondaryCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer secondary.Close()

	client := New(
		WithRetry(RetryConfig{
			MaxRetries:   1,
			Delay:        10 * time.Millisecond,
			FailoverURLs: []string{secondary.URL},
			Statuses:     []int{http.StatusInternalServerError},
		}),
	)

	_, err := client.Get(primary.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !primaryCalled {
		t.Error("Primary URL not called")
	}

	if !secondaryCalled {
		t.Error("Secondary URL not called")
	}
}

func TestHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "value" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(
		WithHeaders(map[string]string{"X-Test": "value"}),
	)

	resp, err := client.Get(ts.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Error("Headers not set correctly")
	}
}

func TestPostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New()
	form := url.Values{"key": {"value"}}

	resp, err := client.PostForm(ts.URL, form)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Error("PostForm failed")
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(WithTimeout(50 * time.Millisecond))

	_, err := client.Get(ts.URL)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded error, got: %v", err)
	}
}

func TestTLSConfig(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(
		WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)

	_, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("Unexpected TLS error: %v", err)
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("Invalid URL", func(t *testing.T) {
		client := New()
		_, err := client.Get("http://invalid.url")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("Invalid Body", func(t *testing.T) {
		client := New()
		_, err := client.Post("http://valid.url", make(chan int))
		if err == nil {
			t.Error("Expected marshaling error")
		}
	})
}

func TestContextCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := New(WithContext(ctx))
	cancel()

	_, err := client.Get(ts.URL)
	if !errors.Is(err, context.Canceled) {
		t.Error("Expected context canceled error")
	}
}

func TestLogging(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelWarn, // Captura logs WARN corretamente
	}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	client := New(
		WithLogger(true),
		WithRetry(RetryConfig{
			MaxRetries: 1,
			Delay:      10 * time.Millisecond,
			Statuses:   []int{http.StatusServiceUnavailable},
			EnableLog:  true,
		}),
	)

	client.Logger = logger

	if httpClient, ok := client.ClientHTTP.(*http.Client); ok {
		if transport, ok := httpClient.Transport.(*RetryTransport); ok {
			transport.Logger = logger
		}
	}

	client.Get(ts.URL)

	time.Sleep(50 * time.Millisecond)

	logOutput := logBuffer.String()
	fmt.Println("Captured log:", logOutput) // Depuração

	if !strings.Contains(logOutput, "Retrying RoundTrip request") {
		t.Errorf("Expected retry log not found. Log output: %s", logOutput)
	}
}

func TestGetDefaultClient(t *testing.T) {
	client1 := GetDefaultClient()
	client2 := GetDefaultClient()

	if client1 != client2 {
		t.Error("Expected GetDefaultClient to return the same instance")
	}

	if client1 == nil {
		t.Error("Expected GetDefaultClient to return a valid instance")
	}
}

func TestWithHTTPClientConfig(t *testing.T) {
	// Create a custom HTTP client configuration
	cfg := &HTTPClientConfig{
		Timeout:           5 * time.Second,
		DisableKeepAlives: true,
		MaxIdleConns:      5,
		MaxConnsPerHost:   2,
	}

	// Create a new client with the provided configuration
	client := New(WithHTTPClientConfig(cfg))

	// Check if the client has the expected configuration
	httpClient, ok := client.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatal("Expected ClientHTTP to be of type *http.Client")
	}

	// Verify timeout configuration
	if httpClient.Timeout != cfg.Timeout {
		t.Errorf("Expected timeout %v, got %v", cfg.Timeout, httpClient.Timeout)
	}

	// Verify transport configuration
	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected Transport to be of type *http.Transport")
	}

	if transport.DisableKeepAlives != cfg.DisableKeepAlives {
		t.Errorf("Expected DisableKeepAlives to be %v, got %v", cfg.DisableKeepAlives, transport.DisableKeepAlives)
	}

	if transport.MaxIdleConns != cfg.MaxIdleConns {
		t.Errorf("Expected MaxIdleConns to be %d, got %d", cfg.MaxIdleConns, transport.MaxIdleConns)
	}

	if transport.MaxConnsPerHost != cfg.MaxConnsPerHost {
		t.Errorf("Expected MaxConnsPerHost to be %d, got %d", cfg.MaxConnsPerHost, transport.MaxConnsPerHost)
	}
}

func TestClientLog(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Ensure INFO level is enabled
	}))

	client := New(
		WithLogger(true), // Enable logging
	)
	client.Logger = logger // Set custom logger

	// Call the log function using structured logging
	client.log("Test message", slog.String("key", "log check"))

	// Capture log output
	logOutput := logBuffer.String()
	fmt.Println("Captured log:", logOutput) // Debugging output

	if !strings.Contains(logOutput, "Test message") || !strings.Contains(logOutput, "log check") {
		t.Errorf("Expected log output not found. Log: %s", logOutput)
	}

	// Test when logging is disabled
	client.EnableLogger = false
	logBuffer.Reset() // Clear buffer
	client.log("This should not be logged")

	if logBuffer.Len() != 0 {
		t.Errorf("Log message should not have been written, but got: %s", logBuffer.String())
	}
}

func TestClientMethods(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer ts.Close()

	// Ensure the singleton client is initialized
	GetDefaultClient()

	// Test GET
	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error in Get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for Get, got %d", resp.StatusCode)
	}

	// Test POST
	resp, err = Post(ts.URL, map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Unexpected error in Post: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for Post, got %d", resp.StatusCode)
	}

	// Test PUT
	resp, err = Put(ts.URL, map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Unexpected error in Put: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for Put, got %d", resp.StatusCode)
	}

	// Test DELETE
	resp, err = Delete(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error in Delete: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for Delete, got %d", resp.StatusCode)
	}

	// Test PostForm
	form := url.Values{"key": {"value"}}
	resp, err = PostForm(ts.URL, form)
	if err != nil {
		t.Fatalf("Unexpected error in PostForm: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for PostForm, got %d", resp.StatusCode)
	}
}

func TestClientPutAndDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate method
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"updated"}`))
		} else if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer ts.Close()

	client := New()

	// Test PUT request
	resp, err := client.Put(ts.URL, map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Unexpected error in Put: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for Put, got %d", resp.StatusCode)
	}
	if !strings.Contains(string(resp.Body), "updated") {
		t.Errorf("Expected response body to contain 'updated', got: %s", string(resp.Body))
	}

	// Test DELETE request
	resp, err = client.Delete(ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error in Delete: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 for Delete, got %d", resp.StatusCode)
	}
}

func TestHTTPTransportOptions(t *testing.T) {
	client := New(
		WithDisableKeepAlives(true),
		WithMaxIdleConns(50),
		WithMaxConnsPerHost(20),
		WithMaxIdleConnsPerHost(10),
	)

	// Ensure client has an HTTP transport
	httpClient, ok := client.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatal("Expected ClientHTTP to be of type *http.Client")
	}

	// Ensure transport is properly set
	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected Transport to be of type *http.Transport")
	}

	// Verify `DisableKeepAlives`
	if !transport.DisableKeepAlives {
		t.Errorf("Expected DisableKeepAlives to be true, got %v", transport.DisableKeepAlives)
	}

	// Verify `MaxIdleConns`
	if transport.MaxIdleConns != 50 {
		t.Errorf("Expected MaxIdleConns to be 50, got %d", transport.MaxIdleConns)
	}

	// Verify `MaxConnsPerHost`
	if transport.MaxConnsPerHost != 20 {
		t.Errorf("Expected MaxConnsPerHost to be 20, got %d", transport.MaxConnsPerHost)
	}

	// Verify `MaxIdleConnsPerHost`
	if transport.MaxIdleConnsPerHost != 10 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 10, got %d", transport.MaxIdleConnsPerHost)
	}
}

func TestHTTPTransportOptionsAdvanced(t *testing.T) {
	// Test WithInsecureTLS
	client := New(WithInsecureTLS(true))

	httpClient, ok := client.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatal("Expected ClientHTTP to be of type *http.Client")
	}

	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected Transport to be of type *http.Transport")
	}

	if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Errorf("Expected InsecureSkipVerify to be true, got %v", transport.TLSClientConfig.InsecureSkipVerify)
	}

	// Test WithTransport
	customTransport := &http.Transport{MaxIdleConns: 100}
	client = New(WithTransport(customTransport))

	httpClient, ok = client.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatal("Expected ClientHTTP to be of type *http.Client")
	}

	if httpClient.Transport != customTransport {
		t.Error("Expected custom transport to be set")
	}

	// Test WithCustomHTTPClient
	customHTTPClient := &http.Client{Timeout: 5 * time.Second}
	client = New(WithCustomHTTPClient(customHTTPClient))

	if client.ClientHTTP != customHTTPClient {
		t.Error("Expected custom HTTP client to be set")
	}

	// Test WithTransportConfig
	preConfiguredTransport := &http.Transport{MaxConnsPerHost: 50}
	client = New(WithTransportConfig(preConfiguredTransport))

	httpClient, ok = client.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatal("Expected ClientHTTP to be of type *http.Client")
	}

	if httpClient.Transport != preConfiguredTransport {
		t.Error("Expected pre-configured transport to be set")
	}
}
