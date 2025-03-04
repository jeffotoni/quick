package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"
)

// TestClient_Configurations verifies that all configurations are applied correctly.
func TestClient_Configurations(t *testing.T) {
	// Create a CookieJar to manage cookies automatically.
	jar, _ := cookiejar.New(nil)

	// Create a custom HTTP transport.
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
		MaxIdleConns:        50,
		MaxConnsPerHost:     30,
		MaxIdleConnsPerHost: 10,
		DisableKeepAlives:   false,
	}

	// Create a custom HTTP client.
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar, // Automatically manages cookies.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: customTransport,
	}

	// Create a client with multiple configurations.
	cClient := New(
		WithTimeout(5*time.Second),
		WithDisableKeepAlives(false),
		WithMaxIdleConns(20),
		WithMaxConnsPerHost(20),
		WithMaxIdleConnsPerHost(20),
		WithContext(context.Background()),
		WithCustomHTTPClient(customHTTPClient),
		WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer TEST_TOKEN",
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
		WithRetry(RetryConfig{
			MaxRetries: 2,
			Delay:      1 * time.Second,
			UseBackoff: true,
			Statuses:   []int{500},
			EnableLog:  false,
		}),
	)

	// Verify that the HTTP client has been set correctly.
	if cClient.ClientHTTP != customHTTPClient {
		t.Errorf("Expected custom HTTP client, but it was not set correctly")
	}

	// Verify the timeout setting.
	httpClient, ok := cClient.ClientHTTP.(*http.Client)
	if !ok {
		t.Fatalf("ClientHTTP is not of type *http.Client")
	}
	if httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", httpClient.Timeout)
	}

	// Verify that the CookieJar is assigned.
	if httpClient.Jar != jar {
		t.Errorf("Expected CookieJar to be set, but it was not assigned correctly")
	}

	// Verify the CheckRedirect function behavior.
	if httpClient.CheckRedirect == nil {
		t.Fatalf("Expected CheckRedirect function to be set, but it is nil")
	}

	// Simulate a redirect limit test.
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	err := httpClient.CheckRedirect(req, make([]*http.Request, 3)) // Exceeds limit.
	if err != http.ErrUseLastResponse {
		t.Errorf("Expected CheckRedirect to return http.ErrUseLastResponse, got %v", err)
	}

	// Verify transport settings.
	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport is not of type *http.Transport")
	}
	if transport.MaxIdleConns != 50 {
		t.Errorf("Expected MaxIdleConns 50, got %d", transport.MaxIdleConns)
	}
	if transport.MaxConnsPerHost != 30 {
		t.Errorf("Expected MaxConnsPerHost 30, got %d", transport.MaxConnsPerHost)
	}
	if transport.MaxIdleConnsPerHost != 10 {
		t.Errorf("Expected MaxIdleConnsPerHost 10, got %d", transport.MaxIdleConnsPerHost)
	}
	if transport.DisableKeepAlives != false {
		t.Errorf("Expected DisableKeepAlives false, got %v", transport.DisableKeepAlives)
	}

	// Verify headers.
	if cClient.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got '%s'", cClient.Headers["Content-Type"])
	}
	if cClient.Headers["Authorization"] != "Bearer TEST_TOKEN" {
		t.Errorf("Expected Authorization header to be 'Bearer TEST_TOKEN', got '%s'", cClient.Headers["Authorization"])
	}

	// Verify TLS settings.
	if transport.TLSClientConfig.InsecureSkipVerify != true {
		t.Errorf("Expected InsecureSkipVerify true, got %v", transport.TLSClientConfig.InsecureSkipVerify)
	}
	if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("Expected MinVersion TLS 1.2, got %d", transport.TLSClientConfig.MinVersion)
	}
}
