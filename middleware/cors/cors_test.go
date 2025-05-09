package cors

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var ConfigDefaultTest = Config{
	AllowedOrigins:   []string{"*"}, // Aceita qualquer origem
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
	ExposedHeaders:   []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           600,
	Debug:            false,
}


// successDefaultCorsHeaders defines expected headers for default configuration.
//
// These headers represent what a properly configured CORS response should include
//
// when using ConfigDefaultTest with a request from http://localhost:30
//
// Default header settings for tests 
var successDefaultCorsHeaders = map[string][]string{
	"Access-Control-Allow-Origin":      {"http://localhost:3000"}, // não pode mais ser "*"
	"Access-Control-Allow-Methods":     {"GET, POST, PUT, DELETE, OPTIONS"},
	"Access-Control-Allow-Headers":     {"Origin, Content-Type, Accept, Authorization"},
	"Access-Control-Expose-Headers":    {"Content-Length"},
	"Access-Control-Allow-Credentials": {"true"},
}

var successCustomCorsHeaders = map[string][]string{
	"Access-Control-Allow-Origin":      {"*"},
	"Access-Control-Allow-Methods":     {"GET, POST"},
	"Access-Control-Allow-Headers":     {""}, // No headers allowed
	"Access-Control-Expose-Headers":    {""},
	"Access-Control-Allow-Credentials": {"true"},
}
// testCors provides a framework for testing CORS middleware.
//
// Contains both a handler function and a request configuration
type testCors struct {
	HandlerFunc http.HandlerFunc
	Request     *http.Request
}

// testCorsSuccess is a preconfigured test case for successful CORS validation.
//
// Simulates an OPTIONS preflight request from http://localhost:3000.
//
// Creating a test request to simulate a real request
var testCorsSuccess = testCors{
	HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	},
	Request: func() *http.Request {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "http://localhost:3000") // <--- IMPORTANTE
		req.Header.Set("Access-Control-Request-Method", "POST")
		return req
	}(),
}

// var testCorsSuccess = testCors{
// 	HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	},
// 	Request: httptest.NewRequest(http.MethodOptions, "/", nil),
// }


func isHeaderEqual(got, want []string) bool {
	return reflect.DeepEqual(got, want)
}

// isHeaderEqual compares two header values for equality.
//
// Uses reflect.DeepEqual to handle slice comparison correctly.
//
//  Helper function to check header equality
// Helper function to compare lists of headers
func isHeaderEqualDefault(got, want []string) bool {
	return reflect.DeepEqual(got, want)
}

// go test -v -failfast -count=1 -run ^TestNew$
func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		testCors      testCors
		wantedHeaders map[string][]string
	}{
		{
			name:          "success",
			config:        ConfigDefaultTest,
			testCors:      testCorsSuccess,
			wantedHeaders: successDefaultCorsHeaders,
		},
		{
			name: "success_CustomConfig",
			config: Config{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{}, // vazio
				ExposedHeaders:   []string{},
				MaxAge:           500,
				AllowCredentials: false,
			},
			testCors: testCors{
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
				Request: func() *http.Request {
					req := httptest.NewRequest("OPTIONS", "/", nil)
					req.Header.Set("Origin", "http://another-domain.com")
					req.Header.Set("Access-Control-Request-Method", "POST")
					return req
				}(),
			},
			wantedHeaders: map[string][]string{
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"GET, POST"},
				// Não coloque a linha abaixo quando AllowedHeaders for vazio:
				// "Access-Control-Allow-Headers": {""},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			t.Logf("==== TEST %s ====", tc.name)

			// Create middleware
			h := New(tc.config)
			a := h(tc.testCors.HandlerFunc)

			// Create response recorder
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, tc.testCors.Request)
			resp := rec.Result()

			// Validate headers
			for k, expected := range tc.wantedHeaders {
				got, exists := resp.Header[k]
				if !exists {
					tt.Errorf("Header %s is missing!\nExpected: %v\n", k, expected)
					continue
				}
				if !isHeaderEqual(got, expected) {
					tt.Errorf("Header %s different!\nReceived: %v\nExpected: %v\n", k, got, expected)
				}
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestDefault$
func TestDefault(t *testing.T) {
	t.Run("success_default", func(tt *testing.T) {
		defConfig := Default()

		if defConfig.Debug != ConfigDefaultTest.Debug {
			tt.Errorf("Incorrect debug, expected: %v, received: %v", ConfigDefaultTest.Debug, defConfig.Debug)
		}

		if defConfig.MaxAge != ConfigDefaultTest.MaxAge {
			tt.Errorf("Incorrect MaxAge, expected: %v, received: %v", ConfigDefaultTest.MaxAge, defConfig.MaxAge)
		}

		if !isHeaderEqualDefault(defConfig.AllowedHeaders, ConfigDefaultTest.AllowedHeaders) {
			tt.Errorf("AllowedHeaders different!")
		}

		if !isHeaderEqualDefault(defConfig.AllowedMethods, ConfigDefaultTest.AllowedMethods) {
			tt.Errorf("AllowedMethods different!")
		}

		if !isHeaderEqual(defConfig.AllowedOrigins, ConfigDefaultTest.AllowedOrigins) {
			tt.Errorf("AllowedOrigins different!")
		}
	})
}

// BenchmarkNew measures performance of middleware creation.
//
//  Useful for identifying any initialization bottlenecks.
//
// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(ConfigDefaultTest)
	}
}

// BenchmarkDefault measures performance of default configuration generation.
//
// Helps ensure the Default() function remains efficient.
//
// go test -bench=. -benchtime=1s -benchmem
func BenchmarkDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default()
	}
}
