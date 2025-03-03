package client

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
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
		WithRetry(
			3,         // Maximum number of retries
			"500ms",   // Delay between attempts
			true,      // Use exponential backoff
			"500,502", // HTTP status for retry
			true,      // show Logger
		),
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
		WithRetry(
			3,         // Maximum number of retries
			"500ms",   // Delay between attempts
			true,      // Use exponential backoff
			"503,504", // HTTP status for retry
			true,      // show Logger
		),
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
		WithRetry(
			2,     // Maximum number of retries
			"1s",  // Delay between attempts
			true,  // Use exponential backoff
			"502", // HTTP status for retry
			true,  // show Logger
		),
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
		WithRetry(
			4,     // Maximum number of retries
			"2s",  // Delay between attempts
			true,  // Use exponential backoff
			"504", // HTTP status for retry
			true,  // show Logger
		),
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
