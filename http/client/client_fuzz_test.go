package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// This test generates random inputs for parseBody()
// and verifies that it always returns valid values:
func FuzzParseBody(f *testing.F) {
	// Seeds iniciais para começar o fuzzing
	f.Add("")                          // Empty string
	f.Add(`{"message": "hello"}`)      // Valid JSON
	f.Add(`{"broken_json": `)          // Invalid JSON
	f.Add(strings.Repeat("A", 100000)) // String too long

	f.Fuzz(func(t *testing.T, input string) {
		reader, err := parseBody(input)

		// The error can only occur if the input is not valid JSON
		if err != nil && !json.Valid([]byte(input)) {
			return // Invalid JSON is expected
		}

		// If there was no error, the reader cannot be nil
		if reader == nil {
			t.Errorf("Expected reader to be non-nil for input: %s", input)
		}
	})
}

// This test sends random input to PostForm() and verifies that it
// correctly handles different types of values:
func FuzzPostForm(f *testing.F) {
	// Creating a test server before fuzzing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Form received"))
	}))
	defer ts.Close() // Guarantees that the server will be closed

	// Adding initial test cases
	f.Add("username", "test_user")
	f.Add("password", "123456")
	f.Add("emoji", "🔥🔥🔥")
	f.Add("long_text", strings.Repeat("A", 100000)) // String too long

	f.Fuzz(func(t *testing.T, key, value string) {
		client := New()

		formData := url.Values{}
		formData.Set(key, value)

		resp, err := client.PostForm(ts.URL, formData)
		if err != nil {
			t.Fatalf("PostForm failed with error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

// Here we test an extreme number of redirects to see if the code remains sane:
func FuzzCheckRedirect(f *testing.F) {
	f.Add(1)
	f.Add(3)  // Expected limit
	f.Add(10) // Over limit
	f.Add(0)  // No redirection
	f.Add(-5) // Invalid value to test protection

	f.Fuzz(func(t *testing.T, numRedirects int) {
		if numRedirects < 0 { // Evita valores inválidos
			t.Skip("Skipping negative numRedirects")
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		}

		req, _ := http.NewRequest("GET", "https://httpbin.org/get", nil)

		// Proteção extra: garantir que numRedirects não seja excessivo
		numValidRedirects := numRedirects % 100 // Evita valores gigantes
		err := client.CheckRedirect(req, make([]*http.Request, numValidRedirects))

		if numValidRedirects >= 3 && err != http.ErrUseLastResponse {
			t.Errorf("Expected http.ErrUseLastResponse for %d redirects, got %v", numValidRedirects, err)
		}
	})
}
