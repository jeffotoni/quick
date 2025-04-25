// The BasicAuth middleware implements HTTP Basic Authentication
// to secure specific routes on an HTTP server.
package basicauth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffotoni/quick/internal/concat"
)

// TestBasicAuth is a table-driven test that verifies various authentication scenarios.
//
// The test covers:
//   - Successful authentication with valid credentials
//   - Failed authentication with invalid credentials
//   - Missing credentials case

// Test Methodology:
// 1. Creates middleware with test credentials ("admin", "1234")
// 2. Sets up a test handler that returns 200 OK when authenticated
// 3. Executes multiple test cases with different Authorization headers
// 4. Verifies the response status codes match expectations

// Example Test Cases:
//   Valid credentials:    "Basic YWRtaW46MTIzNA==" (admin:1234)
//   Invalid credentials:  "Basic d3Jvbmd1c2VyOndyb25ncGFzcw==" (wronguser:wrongpass)
//   No credentials:       "" (no Authorization header)

// TestBasicAuth various types of tests
// TestBasicAuth(t *testing.T)
func TestBasicAuth(t *testing.T) {
	// Middleware configuration
	username := "admin"
	password := "1234"
	middleware := BasicAuth(username, password)

	// Creating test handler
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// List of test scenarios
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Authentication successful",
			authHeader:     concat.String("Basic ", base64.StdEncoding.EncodeToString([]byte("admin:1234"))),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid credentials",
			authHeader:     concat.String("Basic ", base64.StdEncoding.EncodeToString([]byte("wronguser:wrongpass"))),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No credentials",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// Running the tests
	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if ti.authHeader != "" {
				req.Header.Set("Authorization", ti.authHeader)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// Check if the response status is as expected
			if rec.Code != ti.expectedStatus {
				tt.Errorf("Test '%s' failed: expected %d, received %d", ti.name, ti.expectedStatus, rec.Code)
			}
		})
	}
}
