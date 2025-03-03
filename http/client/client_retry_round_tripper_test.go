package client

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
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
	cClient := New(
		//WithTransport(&http.Transport{DisableKeepAlives: true}),
		WithRetryRoundTripper(3, "2s", true, "500,502,503", true),
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

// TestWithRetryRoundTripper validates that retry logic works with RoundTripper.
func TestWithRetryRoundTripper_WithHeaders(t *testing.T) {
	// Creating a test server that fails 2 times and then returns success
	ts := httptest.NewServer(testRetryHandler(2, http.StatusInternalServerError, http.StatusOK, `{"message": "Success!"}`))
	defer ts.Close()

	// Criando o cliente com retry configurado
	cClient := New(
		WithTimeout(8*time.Second), // ⬅️ Aumentado para garantir que todas as tentativas sejam feitas
		WithHeaders(map[string]string{"Content-Type": "application/json"}),

		// Habilitando retry com RoundTripper
		WithRetryRoundTripper(
			3,             // Número máximo de tentativas
			"2s",          // Tempo entre tentativas
			true,          // Habilita backoff exponencial
			"500,502,503", // Status para retry
			true,
		),
	)

	// Executando a requisição
	resp, err := cClient.Post(ts.URL, map[string]string{"name": "jeffotoni"})
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	// Verificando o status da resposta
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verificando o body da resposta
	expectedBody := `{"message": "Success!"}`
	if strings.TrimSpace(string(resp.Body)) != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, string(resp.Body))
	}
}
