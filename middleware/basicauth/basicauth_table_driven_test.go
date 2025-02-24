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
