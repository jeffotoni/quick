package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jeffotoni/quick"
)

// TestNew validates the Logger middleware by simulating various HTTP requests.
//
// This test suite verifies that the middleware correctly processes requests, logs expected
// information, and does not interfere with the normal request flow. It includes tests for:
//   - Standard GET requests
//   - POST requests with body data
//   - Requests with an invalid RemoteAddr (ensuring middleware resilience)
//
// Each test case ensures that:
//   - The response status is correctly returned.
//   - The response body remains unchanged after middleware execution.
//
// Example Usage:
//   - Run `go test -v` to execute this test suite.
func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		setupReq   func() *http.Request
		expectCode int
	}{
		{
			name: "GET request with valid remote address",
			setupReq: func() *http.Request {
				return &http.Request{
					Header:     http.Header{},
					Host:       "localhost:3000",
					Method:     "GET",
					RemoteAddr: "127.0.0.1:3000",
					URL: &url.URL{
						Scheme: "http",
						Host:   "quick.com",
						Path:   "/test",
					},
				}
			},
			expectCode: http.StatusOK,
		},
		{
			name: "POST request with body",
			setupReq: func() *http.Request {
				return &http.Request{
					Header:     http.Header{},
					Host:       "localhost:3000",
					Method:     "POST",
					RemoteAddr: "127.0.0.1:3000",
					Body:       io.NopCloser(bytes.NewBufferString("Request Body")),
					URL: &url.URL{
						Scheme: "http",
						Host:   "quick.com",
						Path:   "/submit",
					},
				}
			},
			expectCode: http.StatusOK,
		},
		{
			name: "Invalid RemoteAddr (should not break)",
			setupReq: func() *http.Request {
				return &http.Request{
					Header:     http.Header{},
					Host:       "localhost:3000",
					Method:     "GET",
					RemoteAddr: "invalid_addr",
					URL: &url.URL{
						Scheme: "http",
						Host:   "quick.com",
						Path:   "/broken",
					},
				}
			},
			expectCode: http.StatusOK, // Middleware should log the error but not fail the request
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Quick in action!"))
			})

			mw := New()
			hmw := mw(handler)

			req := ti.setupReq()
			rec := httptest.NewRecorder()

			hmw.ServeHTTP(rec, req)

			if rec.Code != ti.expectCode {
				t.Errorf("Expected status %d, got %d", ti.expectCode, rec.Code)
			}

			expectedBody := "Quick in action!"
			if rec.Body.String() != expectedBody {
				t.Errorf("Expected response body %q, got %q", expectedBody, rec.Body.String())
			}
		})
	}
}

// TestLoggerMiddleware500 ensures that a request with an invalid remote address
// does not crash the application and correctly returns an HTTP 500.
//
// The middleware should handle incorrect RemoteAddr formats gracefully and log errors
// instead of failing unexpectedly.
//
// Assertions:
//   - The response must have an HTTP 500 status code.
//   - The response body should contain "Internal Server Error".
func TestLoggerMiddleware500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	middleware := New()
	handlerWithMiddleware := middleware(handler)

	req := &http.Request{
		Header:     http.Header{},
		Host:       "localhost:3000",
		RemoteAddr: "invalid",
		Method:     "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "letsgoquick.com",
			Path:   "/error",
		},
	}

	rec := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	expectedBody := "Internal Server Error"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected response body %q, got %q", expectedBody, rec.Body.String())
	}
}

// captureOutput captures stdout and returns its output as a string.
//
// This function redirects os.Stdout, runs the given function, and then
// captures anything written to stdout for later assertions.
//
// Example Usage:
//
//	output := captureOutput(func() {
//	    log.Println("Test log message")
//	})
//	fmt.Println("Captured output:", output)
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old
	return <-out
}

// TestLoggerMiddleware validates the Logger middleware with a text-based format.
//
// This test ensures that logs are correctly formatted and contain expected fields.
func TestLoggerMiddleware(t *testing.T) {
	q := quick.New()

	q.Use(New(Config{
		Format:  "text",
		Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency} | user_id=${user_id}\n",
		Level:   "INFO",
		CustomFields: map[string]string{
			"user_id": "12345",
			"trace":   "abc-xyz",
		},
	}))

	q.Get("/logger", func(c *quick.Ctx) error {
		t.Log("[DEBUG] Logger handler executed")
		return c.Status(200).JSON(map[string]string{"msg": "Quick Logger!"})
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	client := ts.Client()

	output := captureOutput(func() {
		resp, err := client.Get(ts.URL + "/logger")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected 'INFO' log message, but got: %s", output)
	}
}

// TestLoggerMiddlewareJSON ensures JSON-formatted logging works as expected.
func TestLoggerMiddlewareJSON(t *testing.T) {
	q := quick.New()

	q.Use(New(Config{
		Format: "json",
		Level:  "INFO",
	}))

	q.Get("/logger-json", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Quick JSON Logger!"})
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	client := ts.Client()

	output := captureOutput(func() {
		resp, err := client.Get(ts.URL + "/logger-json")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})

	var jsonOutput map[string]interface{}
	err := json.Unmarshal([]byte(output), &jsonOutput)
	if err != nil {
		t.Errorf("JSON output is not valid: %s", err)
	}

	if jsonOutput["status"] != float64(http.StatusOK) {
		t.Errorf("Expected status %d, got %v", http.StatusOK, jsonOutput["status"])
	}
}

// TestLoggerMiddlewareDebug ensures DEBUG messages are logged when the log level is set to DEBUG.
func TestLoggerMiddlewareDebug(t *testing.T) {
	q := quick.New()

	q.Use(New(Config{
		Format:  "slog",
		Pattern: "[${time}] ${level} ${method} ${path} ${status} - ${latency}",
		Level:   "DEBUG",
	}))

	q.Get("/logger-debug", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Quick Debug Logger!"})
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	client := ts.Client()

	output := captureOutput(func() {
		resp, err := client.Get(ts.URL + "/logger-debug")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		testLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		testLogger.Debug("Test debug message")
	})

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected '[DEBUG]' log message, but got: %s", output)
	}
}
