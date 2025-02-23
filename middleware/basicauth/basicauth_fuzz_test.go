package basicauth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

// FuzzBasicAuth func to Fuzz
// FuzzBasicAuth(f *testing.F) {
func FuzzBasicAuth(f *testing.F) {
	// Adding initial inputs to the Fuzz
	f.Add("admin", "1234")      // Valid credentials
	f.Add("wronguser", "wrong") // Invalid credentials
	f.Add("", "")               // Empty credentials
	f.Add("admin", "")          // Empty password
	f.Add("", "1234")           // Empty user

	// Running the test with random inputs
	f.Fuzz(func(t *testing.T, username, password string) {
		middleware := BasicAuth("admin", "1234")

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Access allowed"))
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// Generate Authorization header using random credentials
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		req.Header.Set("Authorization", auth)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		// If credentials are "admin:1234", we expect HTTP 200, otherwise 401
		expectedStatus := http.StatusUnauthorized
		if username == "admin" && password == "1234" {
			expectedStatus = http.StatusOK
		}

		// Check if the returned code is what is expected
		if rec.Code != expectedStatus {
			t.Errorf("Para user=%q, pass=%q esperado %d, mas recebeu %d", username, password, expectedStatus, rec.Code)
		}
	})
}
